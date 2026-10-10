# Unit 4 Quiz — Answers

**Q1: B.** A `Cmd` is `func() Msg`. A is a reducer, C isn't a thing, D is a channel (used with `Program.Send` in Unit 10).

**Q2: C.** `nil` = no command. A quitting needs `tea.Quit`; B/D aren't semantics of nil.

**Q3: B.** `Sequence` runs in order; the quit only happens after the save command finished. A runs them concurrently, so the quit may win. C skips terminal restoration. D is nonsense.

**Q4: B.** `Tick` is one-shot; without a new tick nothing fires again.

**Q5: C.** The runtime runs commands off the event loop, which is why they mustn't touch the model.

**Q6:** Errors don't propagate anywhere on their own; the UI would sit in "loading" forever. An `errMsg` lets `Update` show the error and reset state.

**Q7:** Include a request ID (or the query string) in the reply message; `Update` keeps the latest ID and ignores replies whose ID doesn't match. Optionally debounce with a tick.

**Q8:** False. `Tick` starts its timer relative to when it was issued; `Every` aligns to the system clock boundaries.

**Q9:**
```go
type nameMsg struct{ id int; name string }
func fetchName(id int) tea.Cmd {
    return func() tea.Msg {
        n, err := lookup(id)
        if err != nil { return errMsg{err} }
        return nameMsg{id, n}
    }
}
```

**Q10:** It blocks the single-threaded event loop: no rendering, no key handling for 2s. Use `tea.Tick(2*time.Second, …)`.
