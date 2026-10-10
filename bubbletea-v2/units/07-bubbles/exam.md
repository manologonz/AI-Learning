# Unit 7 Exam

**E1 (Apply).** Build a model with a `textinput` and a `viewport`. `enter` appends the input's text to the viewport's content and clears the input. Show fields, `Init`, `Update`, and `View`.

**E2 (Analyze).** Find every bug:
```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    m.input.Update(msg)
    m.spin.Update(msg)
    return m, nil
}
```

**E3 (Analyze).** After adding `spinner`, the app becomes unresponsive to resizes but still typing works; another dev says "forward WindowSizeMsg to the viewport". Explain what they mean and write the resize handling for a viewport with a 2-line header and 1-line footer.

**E4 (Evaluate).** Two designs for keys: raw strings vs `key.Binding` + `help`. Pick one for a 30-shortcut app; justify with three concrete benefits.

**E5 (Apply).** Add a `key.Binding` "delete" (keys `d`, `delete`) and disable it when nothing is selected; show how it disappears from help.

**E6 (Analyze).** Why might `table.Model.SelectedRow()` panic/misbehave on an empty table, and how do you guard it?

**E7 (Apply).** Implement focus cycling among three `textinput`s (tab / shift+tab), making sure only the focused one blinks and receives input.

**E8 (Evaluate).** List vs Table vs Viewport: choose for (a) a file picker with fuzzy filter, (b) a process monitor with columns, (c) a log tail with scroll.

**E9 (Analyze).** In v1 code `help.Model.Width = 80` and `list.DefaultStyles()`; give the v2 equivalents and explain why.

**E10 (Create).** Design the model for a "GitHub issues browser": list on the left, viewport (issue body) on the right, filter box; describe which components, how messages are routed by focus, and how a resize is distributed.
