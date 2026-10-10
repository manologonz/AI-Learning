# Unit 5 Exam — Answers

**E1:**
```go
import tea "charm.land/bubbletea/v2"
p := tea.NewProgram(model{})
func (m model) View() tea.View {
    v := tea.NewView("hi")
    v.AltScreen = true
    v.MouseMode = tea.MouseModeCellMotion
    return v
}
case tea.KeyPressMsg: if msg.String() == "space" { m.on = !m.on }
```
Four changes graded: import path, options→fields, `View` type, `KeyMsg`→`KeyPressMsg`, and `" "`→`"space"`.

**E2:** Alt-screen state is derived from the model on every render; there's no separate "enter" and "exit" you must pair up, and when the program exits the runtime restores the terminal itself.

**E3:**
```go
func (m model) View() tea.View {
    v := tea.NewView(m.body())
    v.AltScreen = m.full
    return v
}
// Update: case "f": m.full = !m.full
```

**E4:** Because omitted = default, focus reporting would be turned **off** after the first frame, so you'd stop getting `FocusMsg`/`BlurMsg`. The flag makes the description differ between frames.

**E5:** (a) inline: result stays in scrollback. (b) alt-screen: full-screen navigation. (c) inline: brief prompt; leaves a normal terminal. (d) alt-screen: needs full canvas (like `less`); or inline with a paged viewport if output should persist.

**E6:**
```go
c := tea.NewCursor(len(prefix)+len(m.text), 0)   // beware: use display width for wide runes
c.Shape = tea.CursorBar
c.Blink = false
v.Cursor = c
```
(Bonus: use `lipgloss.Width` or `ansi.StringWidth` rather than `len` for non-ASCII.)

**E7:** `OnMouse` only *intercepts* messages; the terminal must be told to report mouse events via `MouseMode` (Cell/AllMotion). Without it no mouse events are sent.

**E8:** With a derived field there's a single source of truth. Two components both toggling via commands could enable/disable in conflicting orders (one disables what the other just enabled). When derived, the parent computes `MouseMode` once from combined state (e.g. `if a.wantsMouse || b.wantsMouse`).

**E9:** `View` runs after every `Update` (including each tick/keypress), so a 40ms render caps you near 25 FPS and delays input handling. Mitigations: cache rendered subcomponents that haven't changed, avoid heavy allocation, or reduce update frequency of noisy messages.

**E10:**
```go
v := tea.NewView(m.render())
v.AltScreen = true
v.MouseMode = tea.MouseModeCellMotion
v.WindowTitle = "▶ " + m.track.Title
v.ProgressBar = tea.NewProgressBar(tea.ProgressBarDefault, m.pctPlayed)
v.ReportFocus = true
```
Messages: `tea.MouseClickMsg` (compute seek from X against bar bounds), `tea.BlurMsg` → pause if `m.pauseOnBlur`, `tea.FocusMsg` → resume, a `tickMsg` to advance position.
