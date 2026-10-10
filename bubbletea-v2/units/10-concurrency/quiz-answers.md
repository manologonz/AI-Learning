# Unit 10 Quiz — Answers

**Q1: B.** Single loop; that's why no locks are needed. A/C would race; D is before the loop exists.

**Q2: B.** `Send` is goroutine-safe. A races; C breaks the runtime's contract; D corrupts the display.

**Q3: C.** Documented no-op. Not a panic, block, or error.

**Q4: B.** Re-arm from `Update`. A blocks the command from delivering; C doesn't loop; D is false.

**Q5: C.** `Kill` skips the final render and `Run` returns `ErrProgramKilled`. A/B are graceful; D only waits.

**Q6:** Intercept, replace, or drop any message before `Update` (e.g. veto `QuitMsg` when there are unsaved changes).

**Q7:** `Send` blocks behind a busy loop; every message triggers `Update` + `View`. Coalesce/throttle at the producer (batch lines per frame, sample at ~30-60Hz) so the UI receives fewer, larger messages.

**Q8:** False. Command goroutines can't be cancelled by the runtime; they finish on their own (use contexts).

**Q9:**
```go
type lineMsg string
type streamEndedMsg struct{}
func waitForLine(ch <-chan string) tea.Cmd {
    return func() tea.Msg {
        s, ok := <-ch
        if !ok { return streamEndedMsg{} }
        return lineMsg(s)
    }
}
```

**Q10:** Concurrent read while `Update` writes a copy/shared references → data race (reads torn/stale data, race detector fires). Send a message asking for the value, or log inside `Update`.
