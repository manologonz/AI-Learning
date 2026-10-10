// Unit 11: multi-field form with validation + a table of submitted rows.
package main

import (
	"errors"
	"fmt"
	"net/mail"
	"os"

	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var errStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

type model struct {
	inputs  []textinput.Model // 0 = name, 1 = email
	focus   int
	err     error
	table   table.Model
	rows    []table.Row
}

func validateName(s string) error {
	if len(s) < 2 {
		return errors.New("name needs at least 2 characters")
	}
	return nil
}

func validateEmail(s string) error {
	if _, err := mail.ParseAddress(s); err != nil {
		return errors.New("not a valid email")
	}
	return nil
}

func newModel() model {
	name := textinput.New()
	name.Placeholder = "Ada Lovelace"
	name.Validate = validateName
	name.Focus()
	email := textinput.New()
	email.Placeholder = "ada@example.com"
	email.Validate = validateEmail

	t := table.New(
		table.WithColumns([]table.Column{{Title: "Name", Width: 20}, {Title: "Email", Width: 28}}),
		table.WithHeight(6),
	)
	return model{inputs: []textinput.Model{name, email}, table: t}
}

func (m model) Init() tea.Cmd { return textinput.Blink }

func (m model) setFocus(i int) (model, tea.Cmd) {
	m.inputs[m.focus].Blur()
	m.focus = (i + len(m.inputs)) % len(m.inputs)
	return m, m.inputs[m.focus].Focus()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyPressMsg); ok {
		switch k.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "tab", "down":
			return m.setFocus(m.focus + 1)
		case "shift+tab", "up":
			return m.setFocus(m.focus - 1)
		case "enter":
			for _, in := range m.inputs {
				if err := in.Err; err != nil || in.Value() == "" {
					m.err = errors.New("fix the form before submitting")
					return m, nil
				}
			}
			m.rows = append(m.rows, table.Row{m.inputs[0].Value(), m.inputs[1].Value()})
			m.table.SetRows(m.rows)
			for i := range m.inputs {
				m.inputs[i].SetValue("")
			}
			m.err = nil
			return m.setFocus(0)
		}
	}
	// Forward to the focused input only.
	var cmd tea.Cmd
	m.inputs[m.focus], cmd = m.inputs[m.focus].Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	s := "Name:  " + m.inputs[0].View() + "\nEmail: " + m.inputs[1].View() + "\n"
	if e := m.inputs[m.focus].Err; e != nil {
		s += errStyle.Render(e.Error()) + "\n"
	}
	if m.err != nil {
		s += errStyle.Render(m.err.Error()) + "\n"
	}
	s += "\n" + m.table.View() + "\n\ntab: next • enter: submit • ctrl+c: quit\n"
	return tea.NewView(s)
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
