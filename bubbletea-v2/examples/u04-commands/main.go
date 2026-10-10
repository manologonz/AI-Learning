// Unit 4: commands - ticks, HTTP, Batch, Sequence, error messages.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
)

type tickMsg time.Time
type statusMsg struct{ code int }
type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

// A Cmd is just func() tea.Msg. This one closes over `url`.
func checkURL(url string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return errMsg{err}
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return errMsg{err} // errors become messages
		}
		defer res.Body.Close()
		return statusMsg{res.StatusCode}
	}
}

// tea.Tick fires ONCE; re-issue it from Update to loop.
func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}

type model struct {
	secs    int
	loading bool
	status  string
}

func (m model) Init() tea.Cmd { return tick() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.secs++
		return m, tick()
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.loading = true
			m.status = ""
			return m, checkURL("https://example.com")
		}
	case statusMsg:
		m.loading = false
		m.status = fmt.Sprintf("HTTP %d", msg.code)
	case errMsg:
		m.loading = false
		m.status = "error: " + msg.Error()
	}
	return m, nil
}

func (m model) View() tea.View {
	s := fmt.Sprintf("Uptime: %ds\n\n", m.secs)
	switch {
	case m.loading:
		s += "Checking...\n"
	case m.status != "":
		s += m.status + "\n"
	}
	s += "\nenter: check example.com • q: quit\n"
	return tea.NewView(s)
}

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
