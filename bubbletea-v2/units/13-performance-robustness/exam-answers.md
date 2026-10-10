# Unit 13 Exam — Answers

**E1:** Keep `lines []string` (or a ring buffer) and a `viewport`; render only `height` lines from `offset` (Bubbles `viewport` with `SetContent` updated only on new data, or manual windowing). Pre-wrap/format on arrival, not per frame; cache the rendered window until offset/size/data changes; append via commands/`Send` batching; avoid re-rendering the full string for a typing box: only the input changes.

**E2:**
```go
go func() {
    defer func() { if r := recover(); r != nil { p.Send(errMsg{fmt.Errorf("worker panic: %v", r)}) } }()
    p.Send(resultMsg(work()))
}()
```

**E3:** (a) fatal at startup: print to stderr with the path and how to create it; exit 1 *before* starting the TUI. (b) transient status-bar error with retry key, keep old data visible. (c) inline validation message at the field; no exit. (d) prominent error, keep unsaved data in memory, offer retry/save-as, don't quit.

**E4:** `case "e": return m, tea.ExecProcess(exec.Command(ed, m.path), func(err error) tea.Msg { return editorDoneMsg{err} })`; `case editorDoneMsg: if msg.err != nil { m.err = msg.err; return m, nil }; m.st = stLoading; return m, loadFile(m.path)`; `case fileLoadedMsg: …`. Note editors exit non-zero on `:cq`; treat as "cancelled", not necessarily fatal.

**E5:** Styles are cheap-ish builders but 200 per frame at 60 FPS adds allocations/GC pressure. Decide by benchmarking `View` (`go test -bench -benchmem`) and profiling; if it shows up, define styles once and derive with `.Width()`.

**E6:** Rarely: when you want a crash to propagate (core dumps, custom crash reporter, or tests that expect the panic). Accept the risk of a raw terminal; ensure your handler restores it (`reset`).

**E7:**
```go
func (m model) View() tea.View {
    if m.w < 40 || m.h < 10 { return tea.NewView("Terminal too small (need 40×10)") }
    …
}
```
State is untouched, so when `WindowSizeMsg` grows again the normal view returns. Keep `AltScreen` consistent (Unit 5).

**E8:** Terminal state changed under you (the child wrote to the screen, or a resize happened while away). Response: after resume, request the size (`tea.RequestWindowSize`) and force a redraw (`tea.ClearScreen` command) if needed; handle `WindowSizeMsg` normally.

**E9:** Order: (1) on quit request, if dirty → confirm (filter/veto, Unit 10); (2) cancel the download context; (3) `tea.Sequence(saveCmd, closeDBCmd, tea.Quit)` or handle `savedMsg` → `closeDB` → `tea.Quit` so failures can stop the exit; (4) `defer` cleanup after `Run` returns for resources the model doesn't own.

**E10:** e.g. (1) no blocking I/O in Update/View; (2) all commands have timeouts/contexts and error messages; (3) guard zero/tiny window sizes; (4) all keys have `"space"`/`KeyPressMsg` v2 forms; (5) bounds checks on cursors/indices for empty data; (6) goroutines have recover + context; (7) run with `-race` and a panic-injection test; (8) works with no TTY (graceful message), `NO_COLOR`, light and dark themes, and small terminals.
