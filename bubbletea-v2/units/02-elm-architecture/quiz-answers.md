# Unit 2 Quiz — Answers

**Q1: B.** v2's exact interface. A has the wrong `Init` (v2 returns only a `Cmd`) and the old v1 `View() string`. C drops the command from `Update`. D has no return values, so Update couldn't produce a new model.

**Q2: C.** The runtime renders after every `Update`. A/B/D: your code describes; it never invokes drawing.

**Q3: B.** You changed a copy, then returned an unrelated fresh model, so the runtime stores the zero model. A assumes mutation persists without returning. C: compiles fine. D: nothing panics.

**Q4: B.** A `Cmd` (`func() Msg`) returning `QuitMsg`; the runtime shuts down and restores the terminal. A is dangerous (leaves terminal raw). C/D wrong kinds.

**Q5: B.** v2 returns `"space"` (v1 returned `" "`). `Key().Text` is still `" "`. C/D are invented strings.

**Q6:** `Update` returns the modified copy and the runtime keeps *that* returned value as the new model, so persistence comes from the return value, not from mutation of the caller's original.

**Q7:** The final `tea.Model` and an `error`.

**Q8:** Advantage: no reallocation, simpler code. Risk: maps are references; old and new model copies share it, so previous-state comparisons/undo history break and concurrent access (e.g. from a goroutine) races.

**Q9:** The renderer owns the terminal; raw prints corrupt the layout and get overwritten. Log to a file with `tea.LogToFile` (Unit 13).

**Q10:**
```go
case "p":
    m.paused = !m.paused
case "q":
    return m, tea.Quit
```
(followed by the shared `return m, nil` at the end of `Update`).
