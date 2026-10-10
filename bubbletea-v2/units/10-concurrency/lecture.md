# Unit 10: Concurrency & External Events

> Runnable code: `examples/u10-external/main.go` — a goroutine pushes events via `Program.Send`.

## Objectives
By the end of this unit you can:
- Explain Bubble Tea's concurrency model (single-threaded `Update`, concurrent commands).
- Push events from outside the program with `Program.Send`, safely from any goroutine.
- Build long-lived subscriptions (channel listeners) with the "re-arm a waiting command" pattern.
- Cancel work with `context`, and shut down the program deterministically (`Quit`, `Kill`, `WithContext`, `Wait`).
- Use `tea.WithFilter` to intercept messages, and know what races you must still avoid.

## Prerequisites
Units 4 and 8, Go goroutines/channels basics.

## Core content

### The concurrency model in three sentences
1. **`Update` and `View` run on one goroutine**, one message at a time, so model state needs no locks.
2. **Commands run concurrently**, each in its own goroutine spawned by the runtime; the message a command returns is sent back into the loop.
3. **Anything outside** (your own goroutines, other packages) talks to the app by **sending messages**, never by touching the model.

A note from the source: command goroutines are *not cancellable* by the runtime. On shutdown a slow command keeps running until it returns, so give your own commands contexts/timeouts.

### Program.Send: the door from the outside
```go
p := tea.NewProgram(model{})
go func() {
    for line := range logs {              // logs is a chan string fed by some other package
        p.Send(logLineMsg(line))
    }
}()
_, err := p.Run()
```
`Send` is goroutine-safe. Semantics (from the implementation): if the program hasn't started, `Send` **blocks** until it does; after the program has exited it's a **no-op** (safe to call late). Note that means a goroutine sending *before* `Run` starts will wait, which is fine as long as `Run` is called eventually.

Use `Send` when the event source is **outside** the Bubble Tea world: websockets, file watchers, OS signals, a job runner. Use **commands** when the work is *started by* the UI in response to a message.

### Subscriptions with commands (no `Send` needed)
Sometimes you want the model to own a stream. Pattern: a command waits for **one** item from a channel, returns it as a message; `Update` handles it and returns the *same wait command again*:
```go
type activityMsg string

func waitForActivity(ch <-chan string) tea.Cmd {
    return func() tea.Msg { return activityMsg(<-ch) }     // blocks (off-loop) until an item arrives
}

func (m model) Init() tea.Cmd { return waitForActivity(m.events) }

case activityMsg:
    m.log = append(m.log, string(msg))
    return m, waitForActivity(m.events)      // re-arm: without this you'd get exactly ONE event
```
This is the same re-arm idea as `Tick`. Closing the channel: `<-ch` returns `""` forever, so use `msg, ok := <-ch` and return a `doneMsg` when `!ok`, else you spin in a tight loop.

### Cancellation and lifetimes
Give background work a `context.Context` derived from the program's lifetime:
```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()                               // stops producers when main returns
p := tea.NewProgram(m, tea.WithContext(ctx)) // cancelling ctx ends the program (Run returns ErrProgramKilled)
```
Inside long commands, `select` on `ctx.Done()` so quitting doesn't leave zombies. For per-request cancel (e.g. superseded searches) store a `cancel` func in the model, call it before starting a new request, and pass the new ctx into the command (copy the ctx into a local; don't read `m.ctx` from the goroutine).

### Shutting down
| Call | Effect |
|---|---|
| `return m, tea.Quit` | graceful, from inside `Update` |
| `p.Quit()` | graceful, from outside (sends `QuitMsg`); no-op if not running |
| `p.Kill()` | immediate; skips final render; `Run` returns `ErrProgramKilled` |
| `p.Wait()` | block until the program has fully shut down |
| `WithContext(ctx)` | cancel context → program exits |
| `tea.Interrupt` / `InterruptMsg` | returns `ErrInterrupted` from `Run` (Ctrl-C style abort) |

Sequence for clean exit: stop producers (cancel ctx) → `p.Quit()` → `p.Wait()`/return from `Run`.

### `tea.WithFilter`
A program-level hook that sees every message *before* `Update` and may replace or drop it:
```go
func filter(m tea.Model, msg tea.Msg) tea.Msg {
    if _, ok := msg.(tea.QuitMsg); ok && m.(model).unsaved { return nil }   // veto quit
    return msg
}
p := tea.NewProgram(model{}, tea.WithFilter(filter))
```
Handy for "unsaved changes" guards or global message rewriting; keep it tiny (it runs on the event loop).

### Races that still exist
- Sharing a **map/slice/pointer** between model and a command or goroutine (Unit 4 rule).
- Reading the model from a goroutine started elsewhere ("just checking `m.count`") — a data race even if it "works".
- Sending on a channel that `Update` also receives from synchronously.
Run `go test -race` and `go run -race .` regularly.

### Backpressure
`Send` blocks when the loop is busy. A firehose producer (10k msgs/s) will stall behind slow `Update`/`View`. Coalesce/throttle at the producer (batch lines into one message, or sample at a fixed rate) instead of forwarding every event.

## Worked examples

### Straightforward: a background ticker via Send
See `examples/u10-external`: a goroutine + `time.Ticker` calls `p.Send(logLineMsg(...))` every 500ms; `Update` appends to a bounded slice.

### Typical: channel-fed activity monitor with re-arming wait
Model owns `events <-chan string`; `Init` returns `waitForActivity`; each `activityMsg` is handled and re-armed. A `doneMsg` on channel close stops re-arming and shows "stream ended".

### Edge case: leaked goroutines and reading the model from outside
```go
go func() {
    for {
        time.Sleep(time.Second)
        p.Send(tickMsg{})              // never stops, even after the program exits
    }
}()
```
After `Run` returns, `Send` is a no-op but this goroutine keeps waking forever (a leak in long-lived processes such as SSH servers running one program per session, Unit 14). Always tie producers to a `context` cancelled when `Run` returns. And the subtle one: a debug goroutine printing `m.count` reads the model concurrently with `Update`'s copies: a data race the race detector will flag.

## Common misconceptions
- **"I need a mutex around the model."** Not if only `Update`/`View` touch it.
- **"`p.Send` is only for tests."** It's the standard bridge to external event sources.
- **"A `Cmd` listening on a channel keeps listening."** It returns one message; re-arm it.
- **"Quitting cancels commands."** They run to completion; use contexts.
- **"`Send` returns an error."** It doesn't; it blocks or no-ops.
- **"More goroutines = faster UI."** The bottleneck is `Update`/`View`; coalesce events.

## Check your understanding
1. Which goroutine runs `Update`? *(The single event-loop goroutine.)*
2. What happens if you call `p.Send` after `Run` returned? *(No-op.)*
3. Why re-arm a channel-wait command? *(Each returns one message.)*
4. How do you veto quitting with unsaved changes? *(`WithFilter` returning nil for `QuitMsg`.)*

## Summary
- One goroutine for `Update`/`View`; commands and external goroutines feed it messages.
- `Program.Send` bridges outside events; re-arming commands build subscriptions.
- Use contexts to cancel; runtime doesn't cancel commands for you.
- Avoid sharing mutable state; throttle high-frequency producers.
