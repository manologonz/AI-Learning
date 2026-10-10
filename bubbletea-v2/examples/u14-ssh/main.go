// Unit 14: serve a Bubble Tea app over SSH with Wish, plus TTY detection.
package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/wish/v2"
	wishtea "charm.land/wish/v2/bubbletea"
	"charm.land/ssh"
	"github.com/charmbracelet/x/term"
)

type model struct{ who string }

func (m model) Init() tea.Cmd { return nil }
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyPressMsg); ok && (k.String() == "q" || k.String() == "ctrl+c") {
		return m, tea.Quit
	}
	return m, nil
}
func (m model) View() tea.View { return tea.NewView("Hello over SSH, " + m.who + "! (q quits)\n") }

func main() {
	// Local mode: refuse to run a TUI when there is no terminal (pipes, CI).
	if len(os.Args) > 1 && os.Args[1] == "local" {
		if !term.IsTerminal(os.Stdin.Fd()) {
			fmt.Fprintln(os.Stderr, "needs an interactive terminal")
			os.Exit(2)
		}
		if _, err := tea.NewProgram(model{who: "local"}).Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort("localhost", "23234")),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			// One model per SSH session; MakeOptions wires the session's PTY.
			wishtea.Middleware(func(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
				return model{who: sess.User()}, wishtea.MakeOptions(sess)
			}),
		),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			fmt.Fprintln(os.Stderr, err)
			done <- os.Interrupt
		}
	}()
	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = s.Shutdown(ctx)
}
