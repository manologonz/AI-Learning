# Final Exam — Answer Key

## Part A

**A1:** No. Value receiver mutates a copy. Use `func (c *counter) Inc() { c.n++ }`.

**A2:** It doesn't match: the dynamic type is `*stopMsg` and the case is the value type `stopMsg`; falls to `default` (or nothing). Match the exact type.

**A3:** `Init() Cmd`, `Update(Msg) (Model, Cmd)`, `View() View`. v1: `View() string`; v2 returns `tea.View` (content + terminal-mode fields). `Init` returns only a `Cmd` (in both versions; old docs sometimes show otherwise).

**A4:** The final model and an error; lets you read user's choices after the UI exits and is central to tests.

**A5:** space: `"space"`, `" "`; shift+a: `"A"`, `"A"` (`Keystroke()` = `"shift+a"`); ctrl+c: `"ctrl+c"`, `""`; up: `"up"`, `""`.

**A6:** `KeyMsg` is an interface matching press and release; with `ReportEventTypes`, each tap yields both. Use `KeyPressMsg`.

**A7:** (i) `Sequence(save, tea.Quit)`; (ii) `Batch(fetchA, fetchB)`.

**A8:** `Tick` fires once; re-arm by returning `tick()` from the `tickMsg` handler (guard with a flag so you don't start multiple chains).

**A9:** `tea.NewProgram(m)`; in `View`: `v.AltScreen = true; v.MouseMode = tea.MouseModeCellMotion`.

**A10:** Rendered width **52** = 50 (Width includes border+padding) + 2 (margin 1+1). Content width = 50 − 2 (border) − 4 (padding 2+2) = **44**. (v1 would have excluded the border.)

**A11:** (1) Embed the component in your model; (2) forward with `m.c, cmd = m.c.Update(msg)`; (3) return `cmd` (or `tea.Batch(cmds...)`). Dropping the Cmd stops blink/tick/animation.

**A12:** `PasteStartMsg`, `PasteMsg{Content}`, `PasteEndMsg`. It prevents pasted text from triggering shortcuts and is delivered as one message.

## Part B

**B1:** Keys are forwarded to both children (or the parent's handler moves the cursor for all keys). Fix with a `focus` field: 
```go
case tea.KeyPressMsg:
    if m.focus == focusForm { m.form, cmd = m.form.Update(msg) } else { m.list, cmd = m.list.Update(msg) }
    return m, cmd
default: // non-key messages to both, batching cmds
```
Also `Blur()` the unfocused input.

**B2:** Add `seq int` and a `cancel context.CancelFunc`; on each keystroke cancel the previous, `m.seq++`, capture `seq, q` and a new `ctx` in the closure; the command returns `searchMsg{seq, res}` or `searchErrMsg{seq, err}` (ignore `context.Canceled`); in `Update` ignore messages with `msg.seq != m.seq`. Optional debounce with `Tick` + seq.

**B3:** Concurrent read of model state while `Update` runs = data race. Designs: (1) `Update` publishes a snapshot into an `atomic.Pointer`/mutex-protected copy that the goroutine reads; (2) the goroutine sends a `queryMsg{reply chan int}` via `p.Send` and `Update` answers on the channel.

**B4:** Selection stored by index. Store the ID and restore the cursor by ID after refresh (nearest index fallback). Booleans allow impossible combinations (loading && err). Use one `state` enum with a `View` switch.

**B5:** (a) `m2,_ := m.Update(press("a"))` assert `m2.screen == formScreen`; type via successive `KeyPressMsg{Text: "h"}` etc; enter → returns a command producing `formSubmittedMsg`; run the command, feed the message back into `Update`, assert list contains "hi". (b) `tm := teatest.NewTestModel(t, newModel(fakeStore), teatest.WithInitialTermSize(80,24)); tm.Send(...)`/`tm.Type("hi")`; `WaitFor` output contains "hi" with a timeout; send `q`; `FinalModel(t, teatest.WithFinalTimeout(2*time.Second))`. Fake: the store (in-memory), the clock. Avoid hangs by timeouts and always quitting.

**B6:** Causes: `View` rebuilds everything each message (cache/window it); expensive formatting/Lip Gloss on giant strings (precompute, render only visible rows); message flood (`AllMotion`, fast ticks, `Send` firehose): throttle/coalesce. A panic in *your* goroutine isn't recovered by Bubble Tea: the process crashes and the terminal can be left raw; wrap with `defer recover()` and `p.Send(errMsg)`.

**B7:** `core` (no Bubble Tea), `storage`, `ui`, `cmd/tool` (flags, TTY detect, wiring), `server` (Wish). `--json` uses core only. SSH: model created per session inside the Wish handler; shared state outside with synchronization or `Send`; per-session goroutines use `sess.Context()`. Security: validate all input (paths/commands), authenticate/authorize, run unprivileged, cap resources.

**B8:** Guard first: `if w == 0 → empty view`; `if w < min || h < min → "too small" view`; **always set the same view fields** (e.g. `AltScreen`) in every branch so state doesn't flip; keep model state untouched so when the size grows the UI returns exactly as it was.

## Part C
Grade with the Unit 15 rubric. Minimum passing: items 1, 2, 3, 6 fully; 4, 5, 7 substantially. Excellent submissions show core-without-Bubble-Tea, no blocking in `Update`, correct v2 idioms (`KeyPressMsg`, `"space"`, `tea.View`, v2 Width rules), and tests that fail if you break a behavior.
