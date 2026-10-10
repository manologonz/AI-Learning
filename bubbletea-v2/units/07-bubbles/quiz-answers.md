# Unit 7 Quiz — Answers

**Q1: B.** A is v1; C/D don't exist.

**Q2: B.** Components render strings for you to embed in your own `tea.View`.

**Q3: B.** The model is returned by value, so assign it back with the cmd. A discards it, C mismatches return types, D isn't a call.

**Q4: B.** Blinking cursors, spinners, and other animations are driven by returned commands. A is wrong: functionality degrades. C/D don't occur.

**Q5: B.** Setters/getters replaced fields in v2. A is v1. C isn't the pattern. D is invented.

**Q6:** The spinner animates by scheduling its next tick as a command; the tick returns as `spinner.TickMsg` to the *parent's* `Update`. If you only forward messages to the focused child, the spinner stops receiving ticks and freezes.

**Q7:** The `q` typed in the filter box triggers the global quit. Fix: while the child is capturing text (`m.list.FilterState() == list.Filtering`) forward keys to it before/instead of global bindings.

**Q8: D.** Binding doesn't fix terminal limitations (e.g. ctrl+i vs tab). A/B/C are true.

**Q9:**
```go
func (t task) FilterValue() string { return t.name }
// list.Item only needs FilterValue; add Title/Description for the default delegate:
func (t task) Title() string       { return t.name }
func (t task) Description() string { if t.done { return "done" }; return "todo" }
if it, ok := m.list.SelectedItem().(task); ok { /* use it */ }
```

**Q10:** v2 removed hidden global background detection; the app must tell components whether the background is dark or light (from `BackgroundColorMsg`), so styles pick readable colors.
