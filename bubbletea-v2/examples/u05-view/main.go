// Unit 5: declarative terminal features on tea.View.
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	alt, mouse, focus bool
	title           string
	progress        int
	x, y            int
	focused         bool
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "a":
			m.alt = !m.alt
		case "m":
			m.mouse = !m.mouse
		case "f":
			m.focus = !m.focus
		case "+":
			m.progress = min(m.progress+10, 100)
		case "-":
			m.progress = max(m.progress-10, 0)
		}
	case tea.FocusMsg:
		m.focused = true
	case tea.BlurMsg:
		m.focused = false
	case tea.MouseMotionMsg:
		m.x, m.y = msg.X, msg.Y
	}
	return m, nil
}

func (m model) View() tea.View {
	body := fmt.Sprintf(
		"a: alt=%v  m: mouse=%v  f: focus-reports=%v (focused=%v)\n+/-: progress=%d\nmouse at (%d,%d)\nq quits\n",
		m.alt, m.mouse, m.focus, m.focused, m.progress, m.x, m.y)
	v := tea.NewView(body)
	// The declarative part: describe what you WANT; Bubble Tea makes it so.
	v.AltScreen = m.alt
	v.ReportFocus = m.focus
	v.WindowTitle = fmt.Sprintf("View demo — %d%%", m.progress)
	if m.mouse {
		v.MouseMode = tea.MouseModeAllMotion
	}
	if m.progress > 0 {
		v.ProgressBar = tea.NewProgressBar(tea.ProgressBarDefault, m.progress)
	}
	v.Cursor = tea.NewCursor(0, 3)
	return v
}

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
