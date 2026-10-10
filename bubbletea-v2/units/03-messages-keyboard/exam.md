# Unit 3 Exam

**E1 (Apply).** Write the `Update` fragment for a list cursor supporting up/down/j/k, `g` (jump to top), `G` (bottom), bounded within `len(m.items)`.

**E2 (Analyze).** Predict `msg.String()` and `msg.Key().Text` for: (a) `a`, (b) shift+a, (c) enter, (d) ctrl+c, (e) space.

**E3 (Analyze).** Find the bug:
```go
case tea.KeyMsg:
    if msg.String() == "enter" { m.submitted++ }
```
on a terminal where release events are enabled.

**E4 (Apply).** Define `tickMsg`, `errMsg` (wrapping an error) and show a type switch handling both plus `tea.WindowSizeMsg`.

**E5 (Evaluate).** Choose `String()` vs `Key().Text` vs field matching for: (a) a "quit" shortcut, (b) a search box, (c) distinguishing shift+enter from enter. Justify each.

**E6 (Analyze).** Why does the following text editor corrupt "héllo" on backspace, and how do you fix it without a library?
```go
m.text = m.text[:len(m.text)-1]
```

**E7 (Apply).** Write `Update` so that pressing `?` toggles `m.showHelp`, and any other key while help is shown closes it (without doing its normal action).

**E8 (Evaluate).** A design doc says "bind Ctrl+I to 'indent' and Tab to 'next field'". Argue whether this is safe.

**E9 (Create).** Design a message vocabulary (type names + payloads) for a chat client's model: connection status, incoming message, send failure. Write the type declarations.

**E10 (Analyze).** Explain what breaks in this View: `w := m.w/2; strings.Repeat("─", w-2)` on first render, and give a two-line fix.
