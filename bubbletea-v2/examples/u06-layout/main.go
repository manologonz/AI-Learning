// Unit 6: Lip Gloss v2 styling and responsive layout.
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	title = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).Padding(0, 1)
	pane = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).Padding(0, 1)
	dim = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type model struct{ w, h int }

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
	case tea.KeyPressMsg:
		if s := msg.String(); s == "q" || s == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	if m.w == 0 {
		return tea.NewView("") // before the first WindowSizeMsg
	}
	// In v2, Width() INCLUDES border and padding (margin is extra). To size the
	// text *inside* a pane, subtract the style's frame with GetHorizontalFrameSize.
	half := m.w / 2
	inner := half - pane.GetHorizontalFrameSize()
	left := pane.Width(half).Render(lipgloss.NewStyle().Width(inner).Render("Left pane\n" + dim.Render("wraps to fit")))
	right := pane.Width(m.w - half).Render("Right pane")
	row := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	page := lipgloss.JoinVertical(lipgloss.Left, title.Render("Layout demo"), row, dim.Render("q to quit"))
	v := tea.NewView(lipgloss.Place(m.w, m.h, lipgloss.Center, lipgloss.Center, page))
	v.AltScreen = true
	return v
}

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
