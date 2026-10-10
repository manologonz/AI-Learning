// Unit 7: Bubbles components - textinput, spinner, list, help + key bindings.
package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type item struct{ title, desc string }

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type keyMap struct{ Tab, Quit key.Binding }

func (k keyMap) ShortHelp() []key.Binding  { return []key.Binding{k.Tab, k.Quit} }
func (k keyMap) FullHelp() [][]key.Binding { return [][]key.Binding{{k.Tab, k.Quit}} }

type model struct {
	input   textinput.Model
	spin    spinner.Model
	list    list.Model
	help    help.Model
	keys    keyMap
	focusIn bool
	chosen  string
}

func newModel() model {
	ti := textinput.New()
	ti.Placeholder = "your name"
	ti.Focus()
	items := []list.Item{
		item{"Go", "compiled, simple"}, item{"Rust", "safe, fast"}, item{"Zig", "small, explicit"},
	}
	l := list.New(items, list.NewDefaultDelegate(), 40, 12)
	l.Title = "Pick one"
	return model{
		input: ti, list: l, focusIn: true, help: help.New(),
		spin: spinner.New(spinner.WithSpinner(spinner.Dot)),
		keys: keyMap{
			Tab:  key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "switch focus")),
			Quit: key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
		},
	}
}

func (m model) Init() tea.Cmd { return tea.Batch(textinput.Blink, m.spin.Tick) }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Tab):
			m.focusIn = !m.focusIn
			if m.focusIn {
				cmds = append(cmds, m.input.Focus())
			} else {
				m.input.Blur()
			}
			return m, tea.Batch(cmds...)
		case msg.String() == "enter" && !m.focusIn:
			if it, ok := m.list.SelectedItem().(item); ok {
				m.chosen = it.title
			}
		}
	case spinner.TickMsg: // the spinner needs its own ticks forwarded
		var c tea.Cmd
		m.spin, c = m.spin.Update(msg)
		return m, c
	}
	// Forward the message to whichever child has focus.
	var c tea.Cmd
	if m.focusIn {
		m.input, c = m.input.Update(msg)
	} else {
		m.list, c = m.list.Update(msg)
	}
	return m, c
}

func (m model) View() tea.View {
	s := lipgloss.JoinVertical(lipgloss.Left,
		m.spin.View()+" Hello "+m.input.Value(),
		m.input.View(),
		m.list.View(),
		fmt.Sprintf("chosen: %q", m.chosen),
		m.help.View(m.keys),
	)
	return tea.NewView(s)
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
