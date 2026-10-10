// Unit 13: robustness - logging, exec'ing an external program, suspend, quit hygiene.
package main

import (
	"fmt"
	"os"
	"os/exec"

	tea "charm.land/bubbletea/v2"
)

type editorFinishedMsg struct{ err error }

type model struct {
	status string
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "e":
			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "vi"
			}
			// ExecProcess hands the terminal to the child, then resumes us.
			c := exec.Command(editor, "/tmp/bubbletea-course-note.txt")
			return m, tea.ExecProcess(c, func(err error) tea.Msg { return editorFinishedMsg{err} })
		case "ctrl+z":
			return m, tea.Suspend
		}
	case editorFinishedMsg:
		if msg.err != nil {
			m.status = "editor error: " + msg.err.Error()
		} else {
			m.status = "editor closed cleanly"
		}
	case tea.ResumeMsg:
		m.status = "resumed after suspend"
	}
	return m, nil
}

func (m model) View() tea.View {
	return tea.NewView(fmt.Sprintf("e: open $EDITOR • ctrl+z: suspend • q: quit\n\n%s\n", m.status))
}

func main() {
	// stdout belongs to the TUI, so log to a file.
	f, err := tea.LogToFile("debug.log", "course")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()

	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
