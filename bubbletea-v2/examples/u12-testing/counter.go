// Unit 12: a small model designed to be testable.
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	count int
	done  bool
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up":
			m.count++
		case "down":
			if m.count > 0 { // business rule worth testing: never below zero
				m.count--
			}
		case "q":
			m.done = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() tea.View { return tea.NewView(fmt.Sprintf("count: %d\n", m.count)) }

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
