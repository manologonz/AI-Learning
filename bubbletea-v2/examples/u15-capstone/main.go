// Unit 15: Taskmate starter. Compiles and runs; everything marked TODO is yours to build.
package main

import (
	"fmt"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// ---- domain (move to internal/core when you split packages) ----
type Task struct {
	ID    string
	Title string
	Done  bool
	Notes string
	Due   time.Time
}

// ---- messages ----
type tasksLoadedMsg struct{ tasks []Task }
type loadErrMsg struct{ err error }

// loadTasks is a stand-in for a real Store; TODO: read JSON from os.UserConfigDir()/taskmate/tasks.json.
func loadTasks() tea.Msg {
	time.Sleep(300 * time.Millisecond) // pretend I/O (fine here: we're inside a Cmd, off the event loop)
	return tasksLoadedMsg{tasks: []Task{{ID: "1", Title: "Learn Bubble Tea v2"}, {ID: "2", Title: "Ship Taskmate", Done: false}}}
}

// ---- data state (Unit 11) ----
type dataState int

const (
	stLoading dataState = iota
	stError
	stEmpty
	stReady
)

type model struct {
	w, h   int
	isDark bool
	state  dataState
	tasks  []Task
	err    error
	cursor int // TODO: replace with bubbles/list, select by ID
}

func newModel() model { return model{state: stLoading} }

func (m model) Init() tea.Cmd { return tea.Batch(loadTasks, tea.RequestBackgroundColor) }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
	case tea.BackgroundColorMsg:
		m.isDark = msg.IsDark()
	case tasksLoadedMsg:
		m.tasks = msg.tasks
		if len(m.tasks) == 0 {
			m.state = stEmpty
		} else {
			m.state = stReady
		}
	case loadErrMsg:
		m.err, m.state = msg.err, stError
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit // TODO: confirm if unsaved; Sequence(save, Quit)
		case "up", "k":
			m.cursor = max(0, m.cursor-1)
		case "down", "j":
			m.cursor = min(len(m.tasks)-1, m.cursor+1)
		case "space":
			if m.state == stReady { // TODO: toggle by ID, then save via a Cmd
				m.tasks[m.cursor].Done = !m.tasks[m.cursor].Done
			}
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	if m.w == 0 {
		return tea.NewView("")
	}
	if m.w < 40 || m.h < 10 {
		return tea.NewView("Terminal too small (need 40×10)")
	}
	var body string
	switch m.state {
	case stLoading:
		body = "Loading tasks…"
	case stError:
		body = fmt.Sprintf("Error: %v\n\npress r to retry", m.err) // TODO: implement retry
	case stEmpty:
		body = "No tasks yet. Press a to add one."
	case stReady:
		for i, t := range m.tasks {
			cur, box := "  ", "[ ]"
			if i == m.cursor {
				cur = "> "
			}
			if t.Done {
				box = "[x]"
			}
			body += fmt.Sprintf("%s%s %s\n", cur, box, t.Title)
		}
	}
	title := lipgloss.NewStyle().Bold(true).Render("Taskmate")
	status := lipgloss.NewStyle().Faint(true).Render(fmt.Sprintf("%d tasks • q quits", len(m.tasks)))
	v := tea.NewView(lipgloss.JoinVertical(lipgloss.Left, title, "", body, "", status))
	v.AltScreen = true
	v.WindowTitle = "Taskmate"
	return v
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
