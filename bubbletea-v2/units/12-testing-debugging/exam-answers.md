# Unit 12 Exam — Answers

**E1:** See `examples/u12-testing/counter_test.go`. Key points: helper `press(s)`; `TestUpdate_DownNeverBelowZero`; `TestUpdate_QuitReturnsQuitCmd` (`cmd().(tea.QuitMsg)`); check `done` on the returned model.

**E2:**
```go
type model struct{ fetch func(string) (string, error) }
func (m model) loadCmd(url string) tea.Cmd { f := m.fetch; return func() tea.Msg { s, err := f(url); if err != nil { return errMsg{err} }; return bodyMsg(s) } }
// tests:
m := model{fetch: func(string) (string, error) { return "ok", nil }}
if m.loadCmd("u")() != bodyMsg("ok") { t.Fatal() }
m = model{fetch: func(string) (string, error) { return "", errors.New("boom") }}
if _, ok := m.loadCmd("u")().(errMsg); !ok { t.Fatal() }
```

**E3:** Nothing quits the program, so `FinalModel` waits, and it has **no default timeout** (it blocks until `go test`'s own 10-minute timeout). Fix: `tm.Send(tea.KeyPressMsg{Code: 'q', Text: "q"})` before `FinalModel` and pass `teatest.WithFinalTimeout(2*time.Second)`; also set `WithInitialTermSize`.

**E4:** (1) Different color profile (force one with `WithColorProfile`); (2) different terminal size (`WithInitialTermSize`); (3) time/animation (freeze/avoid); (4) line endings/paths (CRLF on Windows, use `.gitattributes`); (5) golden not committed.

**E5:** Update-level: menu cursor bounds, `navMsg` routing, `submitMsg` sets name/screen, `ctrl+c` quits from any screen; View-level: each screen contains expected text; teatest: one happy path end-to-end (navigate, type, submit, see "Hello, X!", quit).

**E6:**
```go
teatest.WaitFor(t, tm.Output(), func(b []byte) bool { return bytes.Contains(b, []byte("Loaded 3 items")) },
    teatest.WithDuration(3*time.Second))
```

**E7:** The message is fine; it's the matching: `String()` returns `"space"` not `" "` (v2), so the handler for `" "` never matches.

**E8:** When behavior spans the real loop: command → message → update → view, focus routing, resize handling, exit codes, or bubbles' internal commands (blink, tick); flows a pure `Update` test can't see.

**E9:** Update-level: `m.seq = 2`; send `resultMsg{seq: 1, …}`; assert model unchanged and no command; then send `resultMsg{seq: 2}` and assert the model updated.

**E10:** `logic_test.go` (validators/filters, pure), `update_test.go` (table-driven per screen), `view_test.go` (substring or golden with fixed size/profile), `e2e_test.go` (1–2 teatest flows with timeouts), helpers `press()`, `newTestModel(fakes)`; inject fake storage/clock; no real network/files (use `t.TempDir()`); run `go test -race ./...`.
