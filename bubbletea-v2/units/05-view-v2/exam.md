# Unit 5 Exam

**E1 (Apply).** Port this v1 snippet fully to v2 (imports, signature, options, keys):
```go
p := tea.NewProgram(model{}, tea.WithAltScreen(), tea.WithMouseCellMotion())
func (m model) View() string { return "hi" }
case tea.KeyMsg: if msg.String() == " " { m.on = !m.on }
```

**E2 (Analyze).** Why does v2's design make "forgot to leave alt-screen" bugs impossible in normal use?

**E3 (Apply).** Write a `View` where pressing `f` (toggling `m.full`) switches between inline and full-screen.

**E4 (Analyze).** What happens if `View` sets `ReportFocus = true` on the first render only (via a `first` flag) and false afterwards? Explain in terms of "each View is a full description".

**E5 (Evaluate).** For each app choose inline or alt-screen: (a) a `deploy` progress display whose result should stay visible, (b) a file manager, (c) a yes/no confirmation prompt, (d) a log viewer with search.

**E6 (Apply).** Show a bar-shaped, non-blinking real cursor at column `len(prefix)+len(m.text)`, row 0, in a one-line prompt view.

**E7 (Analyze).** Someone sets `v.OnMouse` expecting to disable mouse handling. What's actually required to receive mouse events at all?

**E8 (Understand).** Why is it valuable that `MouseMode` derives from the model instead of being toggled by commands? Give a scenario with two independent components that both want to change it.

**E9 (Evaluate).** `View` takes 40ms because it builds a huge string every call. Why does that matter in an MVU loop, and what's one mitigation?

**E10 (Create).** Design the `View` fields for a music player: title shows the current track, progress bar shows position, full-screen, mouse enabled for clicking a seek bar, and pause-on-blur behavior (describe the messages you'd handle).
