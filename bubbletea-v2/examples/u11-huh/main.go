// Unit 11: huh forms - both standalone and embedded in a Bubble Tea model.
package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"charm.land/huh/v2"
	tea "charm.land/bubbletea/v2"
)

type model struct {
	form  *huh.Form
	name  string
	color string
	ok    bool
}

func newModel() model {
	m := model{}
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Name").Value(&m.name).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return errors.New("required")
					}
					return nil
				}),
			huh.NewSelect[string]().Title("Color").Options(
				huh.NewOption("Red", "red"), huh.NewOption("Blue", "blue")).Value(&m.color),
			huh.NewConfirm().Title("Proceed?").Value(&m.ok),
		),
	)
	return m
}

func (m model) Init() tea.Cmd { return m.form.Init() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyPressMsg); ok && k.String() == "ctrl+c" {
		return m, tea.Quit
	}
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}
	if m.form.State == huh.StateCompleted {
		return m, tea.Quit
	}
	return m, cmd
}

func (m model) View() tea.View {
	if m.form.State == huh.StateCompleted {
		return tea.NewView(fmt.Sprintf("name=%s color=%s proceed=%v\n", m.name, m.color, m.ok))
	}
	return tea.NewView(m.form.View())
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
