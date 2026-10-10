// Unit 3: keyboard messages, window size, custom messages.
package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type clearMsg struct{} // a custom message: any type can be a Msg

type model struct {
	w, h int
	log  []string
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "c":
			// A command that yields a custom message.
			return m, func() tea.Msg { return clearMsg{} }
		}
		k := msg.Key()
		m.log = append(m.log, fmt.Sprintf("String=%q Code=%q Text=%q Mod=%v", msg.String(), k.Code, k.Text, k.Mod))
		if len(m.log) > 8 {
			m.log = m.log[1:]
		}
	case clearMsg:
		m.log = nil
	}
	return m, nil
}

func (m model) View() tea.View {
	var b strings.Builder
	fmt.Fprintf(&b, "Terminal: %dx%d\nPress keys (c clears, q quits)\n\n", m.w, m.h)
	for _, l := range m.log {
		b.WriteString(l + "\n")
	}
	return tea.NewView(b.String())
}

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
