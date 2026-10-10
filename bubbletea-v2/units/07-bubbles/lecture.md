# Unit 7: Bubbles — Ready-Made Components

> Runnable code: `examples/u07-bubbles/main.go` (textinput + spinner + list + help/key bindings with focus switching).

## Objectives
By the end of this unit you can:
- Explain a Bubble as "a mini Bubble Tea model" and embed one in your model.
- Forward messages to a component and store the returned model + command (the **embed-forward-batch** pattern).
- Use `textinput`, `spinner`, `list`, `viewport`, `table`, `progress`, and `help` at a working level.
- Define `key.Binding`s and a help `KeyMap`, and match keys with `key.Matches`.
- Handle v2-specific API changes (setters like `SetWidth`, `isDark` styles, `KeyPressMsg`).

## Prerequisites
Unit 6.

## Core content

### A Bubble is a model
Every Bubbles component has the shape you already know: `New(...)`, `Update(msg) (Model, tea.Cmd)`, `View() string` (components return **strings**, not `tea.View`; you embed them into your own view), and sometimes `Init`. It is *your* job to route messages to it and place its output.

```go
import (
    "charm.land/bubbles/v2/textinput"
    "charm.land/bubbles/v2/spinner"
    "charm.land/bubbles/v2/list"
    "charm.land/bubbles/v2/viewport"
    "charm.land/bubbles/v2/table"
    "charm.land/bubbles/v2/progress"
    "charm.land/bubbles/v2/help"
    "charm.land/bubbles/v2/key"
)
```

### The embed–forward–batch pattern (learn this cold)
```go
type model struct { input textinput.Model }         // 1. embed

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // (handle your own keys first, maybe return early)
    var cmd tea.Cmd
    m.input, cmd = m.input.Update(msg)               // 2. forward, and STORE the result
    return m, cmd                                    // 3. return the child's command
}
// with several children:
var cmds []tea.Cmd
m.input,  c = m.input.Update(msg);  cmds = append(cmds, c)
m.spin,   c = m.spin.Update(msg);   cmds = append(cmds, c)
return m, tea.Batch(cmds...)
```
Three ways to break it: (1) forgetting to assign back (`m.input.Update(msg)` alone discards the change), (2) dropping the returned `Cmd` (blinking cursors/spinners stop), (3) not forwarding non-key messages (a spinner's `TickMsg`, blink messages) so animation dies.

### textinput
```go
ti := textinput.New()
ti.Placeholder = "your name"
ti.CharLimit = 40
ti.SetWidth(30)               // v2: setter (was field `Width` in v1)
cmd := ti.Focus()             // returns a Cmd (cursor blink) : RETURN IT from Update/Init
v := ti.Value()               // current text
ti.SetValue("")               // reset
ti.Validate = func(s string) error { … }   // runs on change; error in ti.Err
```
A focused input consumes typing; blurred (`ti.Blur()`) ignores it. Only forward keys to the *focused* one.

### spinner
```go
sp := spinner.New(spinner.WithSpinner(spinner.Dot))
// Init:   return m.spin.Tick
// Update: case spinner.TickMsg: m.spin, cmd = m.spin.Update(msg); return m, cmd
// View:   m.spin.View() + " loading…"
```
The spinner animates by returning ticks; you must keep forwarding `spinner.TickMsg`.

### viewport (scrolling text)
```go
vp := viewport.New(viewport.WithWidth(w), viewport.WithHeight(h))   // v2 options
vp.SetContent(longString)
// on WindowSizeMsg: vp.SetWidth(msg.Width); vp.SetHeight(msg.Height - headerFooter)
// forward messages: m.vp, cmd = m.vp.Update(msg)  // arrows/pgup/pgdn/mouse wheel scroll
```

### list
```go
type item struct{ title, desc string }
func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }   // required by list.Item

l := list.New([]list.Item{ item{…}, item{…} }, list.NewDefaultDelegate(), w, h)
l.Title = "Pick one"
sel, ok := l.SelectedItem().(item)          // type assertion (Unit 1!)
```
It ships filtering (`/`), pagination, status bar, and help. Resize with `l.SetSize(w, h)`. Custom rendering: implement `list.ItemDelegate`.

### table
```go
t := table.New(
    table.WithColumns([]table.Column{{Title: "Name", Width: 20}, {Title: "Email", Width: 28}}),
    table.WithRows([]table.Row{{"Ada", "ada@example.com"}}),
    table.WithFocused(true), table.WithHeight(7),
)
row := t.SelectedRow()                       // table.Row is []string
```

### progress
`progress.New()`; drive it with `SetPercent(0.4)` (returns a Cmd for animation) or `ViewAs(0.4)` for a static render. It needs its `progress.FrameMsg` forwarded for animation.

### Key bindings & help
Hard-coded `msg.String() == "q"` scatters shortcuts. `key.Binding` names a *semantic action*, lists its keys, and carries help text; `help` renders it.
```go
type keyMap struct{ Up, Down, Quit key.Binding }
func (k keyMap) ShortHelp() []key.Binding  { return []key.Binding{k.Up, k.Down, k.Quit} }
func (k keyMap) FullHelp() [][]key.Binding { return [][]key.Binding{{k.Up, k.Down}, {k.Quit}} }

keys := keyMap{
    Up:   key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
    Quit: key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
}
// Update: case tea.KeyPressMsg: if key.Matches(msg, keys.Quit) { return m, tea.Quit }
// View:   h.View(keys)   // help.New(); h.SetWidth(w) so it truncates properly
```
Benefit: change a key in one place; disable with `binding.SetEnabled(false)` and it disappears from help *and* stops matching.

### v2 differences to expect
- Many `Width`/`Height` **fields** → `SetWidth()/Width()` methods.
- `DefaultKeyMap` variables → functions (`DefaultKeyMap()`).
- Styles that used to auto-adapt (help, list…) take `isDark` explicitly (`help.DefaultStyles(isDark)`, `list.DefaultStyles(isDark)`); get it from `tea.BackgroundColorMsg`.
- Everything matches `tea.KeyPressMsg` internally; `NewModel` aliases are gone (`New`).

## Worked examples

### Straightforward: a name prompt
Model holds `textinput.Model`; `Init` returns `textinput.Blink`; `Update` forwards messages, and on `enter` reads `m.input.Value()` and quits. `View`: `tea.NewView("Name: " + m.input.View() + "\n")`.

### Typical: input + list with focus switching (see `examples/u07-bubbles`)
`tab` toggles `focusIn`; `Update` forwards to only the focused child; `spinner.TickMsg` is forwarded *always* (so it keeps animating even when the input isn't focused); `key.Matches` handles global keys first.

### Edge case: components swallow your keys
The list's filter mode consumes letters (including `q`). If your global `q` = quit is checked *before* forwarding, you'd quit while typing "quit" in the filter. Check `m.list.FilterState() == list.Filtering` (or `SettingFilter()`) and skip global shortcuts while the child is capturing input. Similarly, blurred vs focused inputs: forgetting `Blur()` on the old input means two fields both react (if you forward to both) or the cursor keeps blinking on the wrong field.

## Common misconceptions
- **"Components return `tea.View`."** They return strings from `View()`; you wrap the final string.
- **"`m.input.Update(msg)` mutates `m.input`."** It returns the updated component; assign it.
- **"I can ignore the returned Cmd."** Then blink/tick/animation stop.
- **"`textinput.Width = 30`."** v2 uses `SetWidth(30)`.
- **"Global `q` always works."** Not while a child is capturing text.
- **"Help text is separate from bindings."** It lives on the binding, so it can't drift from the keys.

## Check your understanding
1. What three steps make up embed-forward-batch? *(Embed, forward+store, return the cmd/batch.)*
2. What breaks if you forget to forward `spinner.TickMsg`? *(The spinner freezes.)*
3. Why use `key.Binding` over string compares? *(One source of truth; help text; enable/disable.)*
4. How do you read a list's selection safely? *(`SelectedItem().(item)` with `ok` check.)*

## Summary
- Bubbles are embeddable mini-models: forward messages, store the result, return the command.
- Learn textinput, spinner, viewport, list, table, progress, help/key.
- v2: setters, `isDark` styles, functions for default keymaps, `KeyPressMsg`.
- Route input only to the focused child; don't let global shortcuts steal typing.
