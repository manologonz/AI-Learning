// Unit 8: parent/child models and screen routing.
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

// Messages that children send UP to the router.
type navMsg struct{ to screen }
type submitMsg struct{ name string }

type screen int

const (
	menuScreen screen = iota
	formScreen
	doneScreen
)

// ---- child: menu ----
type menu struct{ cursor int }

var menuItems = []string{"Enter name", "Quit"}

func (m menu) Update(msg tea.Msg) (menu, tea.Cmd) {
	if k, ok := msg.(tea.KeyPressMsg); ok {
		switch k.String() {
		case "up", "k":
			m.cursor = max(0, m.cursor-1)
		case "down", "j":
			m.cursor = min(len(menuItems)-1, m.cursor+1)
		case "enter":
			if m.cursor == 0 {
				return m, func() tea.Msg { return navMsg{formScreen} }
			}
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m menu) View() string {
	s := "Main menu\n\n"
	for i, it := range menuItems {
		cur := "  "
		if i == m.cursor {
			cur = "> "
		}
		s += cur + it + "\n"
	}
	return s
}

// ---- child: form ----
type form struct{ text string }

func (f form) Update(msg tea.Msg) (form, tea.Cmd) {
	if k, ok := msg.(tea.KeyPressMsg); ok {
		switch k.String() {
		case "enter":
			name := f.text
			return f, func() tea.Msg { return submitMsg{name} }
		case "esc":
			return f, func() tea.Msg { return navMsg{menuScreen} }
		case "backspace":
			if len(f.text) > 0 {
				f.text = f.text[:len(f.text)-1]
			}
		default:
			f.text += k.Key().Text // Text is "" for non-printing keys
		}
	}
	return f, nil
}

func (f form) View() string { return "Your name: " + f.text + "█\n\nenter: submit • esc: back\n" }

// ---- parent / router ----
type model struct {
	current screen
	menu    menu
	form    form
	name    string
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// 1. Global concerns first.
	if k, ok := msg.(tea.KeyPressMsg); ok && k.String() == "ctrl+c" {
		return m, tea.Quit
	}
	// 2. Messages addressed to the router itself.
	switch msg := msg.(type) {
	case navMsg:
		m.current = msg.to
		return m, nil
	case submitMsg:
		m.name, m.current = msg.name, doneScreen
		return m, nil
	}
	// 3. Everything else goes to the ACTIVE child only.
	var cmd tea.Cmd
	switch m.current {
	case menuScreen:
		m.menu, cmd = m.menu.Update(msg)
	case formScreen:
		m.form, cmd = m.form.Update(msg)
	case doneScreen:
		if _, ok := msg.(tea.KeyPressMsg); ok {
			m.current = menuScreen
		}
	}
	return m, cmd
}

func (m model) View() tea.View {
	switch m.current {
	case formScreen:
		return tea.NewView(m.form.View())
	case doneScreen:
		return tea.NewView(fmt.Sprintf("Hello, %s!\n\n(any key)\n", m.name))
	}
	return tea.NewView(m.menu.View())
}

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
