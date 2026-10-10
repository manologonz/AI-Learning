# Unit 10 Exam — Answers

**E1:**
```go
ctx, cancel := context.WithCancel(context.Background())
p := tea.NewProgram(model{})
go func() { defer cancel(); _, _ = p.Run() }()   // or run in main
go tailLog(ctx, func(s string){ p.Send(logMsg(s)) })
// main: <-ctx.Done(); p.Wait()
```
Simplest correct shape: `go tailLog(ctx, ...)`; `_, err := p.Run()`; `cancel()` after Run returns; producers select on `ctx.Done()`.

**E2:** (1) The command mutates `m.counts` (a shared map) on a goroutine while `View`/`Update` read it: data race; (2) mutating shared state from a command violates "return a message"; (3) `m` captured by the closure is a stale copy. Fix: compute locally, return `countsMsg{...}`, and apply in `Update`.

**E3:**
```go
func waitFor(ch <-chan string) tea.Cmd { return func() tea.Msg { s, ok := <-ch; if !ok { return doneMsg{} }; return itemMsg(s) } }
func (m model) Init() tea.Cmd { return waitFor(m.ch) }
case itemMsg: m.items = append(m.items, string(msg)); return m, waitFor(m.ch)
case doneMsg: m.done = true; return m, nil
```

**E4:** (a) `Send` (external source) — or a channel subscription; (b) command (UI-initiated); (c) `Send` from a signal handler goroutine; (d) command with a context.

**E5:** `Send` blocks until the loop starts, so it's safe (it just waits). If `Run` is never called, that goroutine blocks forever: a leak; so ensure `Run` is called or use a context to exit.

**E6:**
```go
func filter(m tea.Model, msg tea.Msg) tea.Msg {
    if _, ok := msg.(tea.QuitMsg); ok && m.(model).dirty { return nil }
    return msg
}
```
Force quit: a key (`ctrl+q`) whose handler clears `dirty` first then returns `tea.Quit`, or use `tea.Interrupt` / `p.Kill()`.

**E7:** `p.Quit()` is graceful (final render, cleanup) and should be preferred if the program is healthy. `p.Kill()` if hung or for hard deadlines (skips final render). Context cancellation (`WithContext`) integrates cleanly with other cancelled work but exits with `ErrProgramKilled`. Usually: signal handler → `cancel()` producers + `p.Quit()`; fall back to `Kill` after a timeout.

**E8:** Producer goroutine reads lines into a buffer; a `time.Ticker(33ms)` flushes: `if len(buf) > 0 { p.Send(linesMsg(buf)); buf = nil }`. One message per flush.

**E9:** Goroutines/timers that outlive the session: producers not cancelled when `Run` returns; `Tick` chains fine (die with program) but custom `time.Ticker` goroutines, channel subscriptions never closed, `context`s never cancelled. Check with `runtime.NumGoroutine()`/pprof.

**E10:**
```go
type model struct{ cancel context.CancelFunc; seq int }
// on keystroke:
if m.cancel != nil { m.cancel() }
ctx, cancel := context.WithCancel(context.Background())
m.cancel, m.seq = cancel, m.seq+1
seq, q := m.seq, m.query
return m, func() tea.Msg { res, err := search(ctx, q); if err != nil { return searchErrMsg{seq, err} }; return searchMsg{seq, res} }
// on reply: if msg.seq != m.seq { return m, nil }  // stale
```
Cancellation frees resources; seq guards the race where a reply was already in flight. Context errors from cancellation must be ignored, not shown.
