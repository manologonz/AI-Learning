// Unit 10: events from outside the program (Program.Send) and cancellation.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
)

type logLineMsg string

type model struct {
	lines []string
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case logLineMsg:
		m.lines = append(m.lines, string(msg))
		if len(m.lines) > 15 {
			m.lines = m.lines[1:]
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	s := "Background producer (q quits):\n\n"
	for _, l := range m.lines {
		s += l + "\n"
	}
	return tea.NewView(s)
}

func main() {
	p := tea.NewProgram(model{})

	// A goroutine OUTSIDE the event loop pushes messages in via p.Send.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		t := time.NewTicker(500 * time.Millisecond)
		defer t.Stop()
		for i := 0; ; i++ {
			select {
			case <-ctx.Done():
				return
			case now := <-t.C:
				p.Send(logLineMsg(fmt.Sprintf("%s event #%d", now.Format("15:04:05"), i)))
			}
		}
	}()

	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
