// Unit 9: mouse clicks, wheel, paste, focus, key release.
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	events []string
}

func (m model) Init() tea.Cmd { return nil }

func (m model) add(s string) model {
	m.events = append(m.events, s)
	if len(m.events) > 10 {
		m.events = m.events[1:]
	}
	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		m = m.add(fmt.Sprintf("press %s repeat=%v", msg.String(), msg.Key().IsRepeat))
	case tea.KeyReleaseMsg:
		m = m.add("release " + msg.String())
	case tea.MouseClickMsg:
		m = m.add(fmt.Sprintf("click %s at %d,%d", msg.String(), msg.X, msg.Y))
	case tea.MouseWheelMsg:
		m = m.add(fmt.Sprintf("wheel %s", msg.String()))
	case tea.PasteMsg:
		m = m.add(fmt.Sprintf("paste %d bytes", len(msg.Content)))
	case tea.FocusMsg:
		m = m.add("terminal focused")
	case tea.BlurMsg:
		m = m.add("terminal blurred")
	}
	return m, nil
}

func (m model) View() tea.View {
	s := "Click, scroll, paste, switch windows, press keys. q quits.\n\n"
	for _, e := range m.events {
		s += e + "\n"
	}
	v := tea.NewView(s)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	v.ReportFocus = true
	v.KeyboardEnhancements.ReportEventTypes = true // ask for release/repeat events
	return v
}

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
