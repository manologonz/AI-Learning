# Unit 12: Testing & Debugging TUIs

> Runnable code: `examples/u12-testing/` (`counter.go` + `counter_test.go`) — run `go test ./examples/u12-testing/`.

## Objectives
By the end of this unit you can:
- Unit-test `Update` as a pure function, including the commands it returns.
- Run a program headlessly with `teatest` (v2 import path), send keys, wait for output, and assert the final model.
- Use golden files to catch view regressions.
- Debug a running TUI with file logging and (optionally) a debugger, and interpret common failure modes.
- Structure code so most logic is testable without a terminal.

## Prerequisites
Units 2–8.

## Core content

### The testing pyramid for TUIs
1. **Pure logic tests** (fast, many): functions on data, validators, `visible()` from Unit 11.
2. **`Update` tests** (fast, many): call `m.Update(msg)` directly, assert the returned model and command.
3. **`View` tests** (fast): call `m.View().Content` and assert substrings or golden files.
4. **Program tests with `teatest`** (slower, few): run the real event loop headlessly for end-to-end flows.

Design pays off here: because `Update` is a function of (model, message), most behavior needs **no terminal**.

### Level 1 & 2: testing Update directly
```go
func TestUpdate_DownNeverBelowZero(t *testing.T) {
    next, cmd := model{}.Update(tea.KeyPressMsg{Code: tea.KeyDown})
    if got := next.(model).count; got != 0 { t.Fatalf("count=%d", got) }
    if cmd != nil { t.Fatal("expected no command") }
}
```
Construct messages by hand: `tea.KeyPressMsg{Code: 'q', Text: "q"}`, `tea.KeyPressMsg{Code: tea.KeyUp}`, `tea.WindowSizeMsg{Width: 80, Height: 24}`, or your custom messages.

**Testing commands**: a `Cmd` is just a function; call it and inspect the message:
```go
_, cmd := model{}.Update(press("q"))
if _, ok := cmd().(tea.QuitMsg); !ok { … }
```
For commands that do I/O, **inject dependencies** (an interface or function field in the model: `fetch func(string) (Result, error)`) so tests can pass a fake, then call `cmd()` and assert on the message. `tea.Batch` returns a `tea.BatchMsg` (a `[]tea.Cmd`) when there are 2+ commands; type-assert and run each.

Type assertions on the returned model: `next.(model)`; consider a helper `func update(m model, msg tea.Msg) (model, tea.Cmd)` in tests.

### Level 3: testing View
```go
out := m.View().Content
if !strings.Contains(out, "count: 2") { … }
```
Lip Gloss output contains ANSI escapes; either assert substrings of plain text or strip escapes (e.g. `ansi.Strip` from `github.com/charmbracelet/x/ansi`). Force a color profile in tests (`tea.WithColorProfile`, or `lipgloss` writers) so output is deterministic.

### Level 4: teatest
```go
import "github.com/charmbracelet/x/exp/teatest/v2"      // note the /v2

tm := teatest.NewTestModel(t, model{}, teatest.WithInitialTermSize(40, 10))
tm.Send(tea.KeyPressMsg{Code: tea.KeyUp})                 // or tm.Type("hello") for text
teatest.WaitFor(t, tm.Output(), func(b []byte) bool {     // poll the output until it matches
    return bytes.Contains(b, []byte("count: 1"))
}, teatest.WithDuration(2*time.Second))
tm.Send(tea.KeyPressMsg{Code: 'q', Text: "q"})
final := tm.FinalModel(t, teatest.WithFinalTimeout(2*time.Second)).(model)
```
Rules:
- **Always set a size** (`WithInitialTermSize`); otherwise layout code sees 0×0.
- **Never `time.Sleep`**; use `WaitFor` with a condition, and always set timeouts so a hang fails fast.
- `FinalModel` waits for the program to finish and has **no default timeout**: your test must cause a quit, and you should pass `teatest.WithFinalTimeout(...)` so a hang fails fast instead of blocking until `go test` times out.
- `Output()` is a stream; `WaitFor` consumes it, so later assertions read from where you left off.

### Golden files
```go
out, _ := io.ReadAll(tm.FinalOutput(t, teatest.WithFinalTimeout(time.Second)))
teatest.RequireEqualOutput(t, out)      // compares to testdata/<TestName>.golden
```
Run `go test -update` to (re)write goldens. Review diffs like code. Goldens are brittle for animated views; freeze time/spinners, fix sizes and color profiles, or use goldens only for stable screens.

### Debugging a live TUI
You can't `fmt.Println`. Options:
1. **File logging**: `f, _ := tea.LogToFile("debug.log", "myapp"); defer f.Close()`, then `log.Printf(...)` anywhere; in another terminal: `tail -f debug.log`.
2. **Dump messages**: log `%T %+v` at the top of `Update` to see the real message stream; invaluable for "why didn't my key match".
3. **A debugger**: run headless (`dlv debug --headless --listen=:2345`) and attach from an editor; or run the program with `WithInput/WithOutput` pointing to files.
4. **Panics**: Bubble Tea catches panics, restores the terminal, prints the stack and `Run` returns an error wrapping `tea.ErrProgramPanic`. Set `TEA_DEBUG=true` to also write `bubbletea-panic-<unix>.log`. If you need the raw crash, use `WithoutCatchPanics()` (terminal may be left raw).
5. **Race detector**: `go test -race ./...` and `go run -race .` for anything with goroutines (Unit 10).

### Common failure diagnoses
| Symptom | Likely cause |
|---|---|
| Key never triggers | matching `" "` not `"space"`, `KeyMsg` vs `KeyPressMsg`, or another component swallowed it |
| UI frozen | blocking I/O in `Update`/`View` |
| Spinner/blink stopped | dropped a returned `Cmd`, or didn't forward the tick message |
| Layout garbled on start | used size before first `WindowSizeMsg` |
| Test hangs | never quit / no timeout / awaited output that never appears |
| Flaky test | real time/randomness; inject clocks/seeds |

## Worked examples

### Straightforward: table-driven Update tests
```go
tests := []struct{ name string; start model; key tea.KeyPressMsg; want int }{
    {"up",   model{0}, tea.KeyPressMsg{Code: tea.KeyUp},   1},
    {"down0", model{0}, tea.KeyPressMsg{Code: tea.KeyDown}, 0},
    {"down2", model{2}, tea.KeyPressMsg{Code: tea.KeyDown}, 1},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        next, _ := tt.start.Update(tt.key)
        if got := next.(model).count; got != tt.want { t.Errorf("got %d want %d", got, tt.want) }
    })
}
```

### Typical: faking I/O
```go
type model struct { fetch func(context.Context, string) (string, error); … }
func (m model) load(url string) tea.Cmd {
    f := m.fetch                                   // copy the dependency for the closure
    return func() tea.Msg { s, err := f(context.Background(), url); if err != nil { return errMsg{err} }; return bodyMsg(s) }
}
// test:
m := model{fetch: func(context.Context, string) (string, error) { return "ok", nil }}
if msg := m.load("x")(); msg != bodyMsg("ok") { t.Fatal(msg) }
```

### Edge case: async + WaitFor pitfalls
A test that sends a key then immediately asserts on `FinalModel` without quitting hangs (timeout). A test that asserts on output *text that includes a spinner frame* fails intermittently. Another trap: `WaitFor` matching too early on stale text from a previous frame; include something unique to the new state (e.g. a counter value) in the condition. Also a subtle one: `tea.Batch(a, b)` in a test: calling `cmd()` returns a `tea.BatchMsg`, not `a`'s message; unwrap and run.

## Common misconceptions
- **"TUIs can't be unit-tested."** `Update` is just a function.
- **"Use `time.Sleep` to wait."** Poll with `WaitFor` and timeouts.
- **"Golden files are always safe."** Animation/colors/sizes make them flaky.
- **"Print debugging works."** Log to a file.
- **"teatest is v1 only."** There is a `/v2` path matching Bubble Tea v2.
- **"Commands aren't testable."** Call them; inject dependencies.

## Check your understanding
1. How do you inspect a command's result in a test? *(Call it; assert on the returned message.)*
2. Why set `WithInitialTermSize`? *(Otherwise size-dependent views see 0×0.)*
3. Where does debug output go? *(A log file via `tea.LogToFile`.)*
4. What does `Batch(a,b).()` return? *(A `tea.BatchMsg`.)*

## Summary
- Most TUI logic is testable as pure functions: `Update`, `View`, validators.
- `teatest` runs the real loop headlessly; use `WaitFor` with timeouts; set a size.
- Inject I/O to test commands.
- Debug via file logs, message dumps, `TEA_DEBUG`, and `-race`.
