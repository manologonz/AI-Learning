# Unit 10 Exam

**E1 (Apply).** Write `main` for an app that tails a log via a goroutine feeding `p.Send`, with clean shutdown: cancel producers, quit, and wait.

**E2 (Analyze).** Find every concurrency bug:
```go
type model struct{ counts map[string]int }
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg.(type) {
    case startMsg:
        return m, func() tea.Msg {
            for i := 0; i < 100; i++ { m.counts["x"]++ }
            return doneMsg{}
        }
    }
    return m, nil
}
func (m model) View() tea.View { return tea.NewView(fmt.Sprint(m.counts)) }
```

**E3 (Apply).** Implement a channel subscription with re-arming and a `doneMsg` on close; show `Init` and `Update`.

**E4 (Evaluate).** For each, choose `Send` or a command: (a) websocket pushes messages; (b) user presses `r` to refresh; (c) SIGHUP should reload config; (d) a build takes 10s after user confirms.

**E5 (Analyze).** Why might `p.Send` in a goroutine started before `p.Run()` be safe, and what would break if `Run` is never called?

**E6 (Apply).** Add a filter that blocks `QuitMsg` while `m.dirty`, and a way for the user to force quit.

**E7 (Evaluate).** Compare `Quit`, `Kill`, and context cancellation for a Ctrl-C from an unrelated signal handler.

**E8 (Apply).** Debounce/coalesce: turn 5,000 log lines/sec into at most 30 messages/sec with a producer-side batcher. Sketch it.

**E9 (Analyze).** A per-session SSH TUI (one `Program` per connection) grows memory over time. What concurrency leak would you look for?

**E10 (Create).** Design cancellation for a search-as-you-type feature: each keystroke starts a search command and cancels the previous one; describe model fields, the command, and stale-reply protection.
