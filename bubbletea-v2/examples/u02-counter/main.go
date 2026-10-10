// Unit 2: the smallest complete Bubble Tea v2 program.
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

// model is ALL of the app's state.
type model struct {
	count int
}

// Init returns an optional first command. nil = "do nothing".
func (m model) Init() tea.Cmd { return nil }

// Update reacts to a message and returns the next model plus an optional command.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			m.count++
		case "down", "j":
			m.count--
		}
	}
	return m, nil
}

// View turns the model into what the terminal should show. It must not mutate m.
func (m model) View() tea.View {
	return tea.NewView(fmt.Sprintf("Count: %d\n\n↑/↓ to change, q to quit\n", m.count))
}

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
