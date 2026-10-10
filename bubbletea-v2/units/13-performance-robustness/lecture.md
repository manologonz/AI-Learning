# Unit 13: Performance, Errors & Robustness

> Runnable code: `examples/u13-exec/main.go` (logging, `ExecProcess`, suspend/resume).

## Objectives
By the end of this unit you can:
- Explain the per-message cost model (`Update` + `View` + render) and find hot spots.
- Keep `Update`/`View` fast: cache, window, avoid allocation storms, throttle noisy messages.
- Handle errors as data: user-visible error states, retry, and never-crash policies.
- Explain what happens on panic and how the terminal is restored; know `ErrProgramPanic`, `TEA_DEBUG`.
- Hand the terminal to another program (`ExecProcess`), suspend/resume, print above the UI, and degrade gracefully when there's no TTY.

## Prerequisites
Units 4, 10, 12.

## Core content

### The cost model
Every message costs: `Update` → `View` → renderer diff → terminal write. The v2 renderer diffs the new frame against the previous cell buffer and writes only changes, capped at a max FPS (default 60; `tea.WithFPS(n)`, max 120). So:
- Slow **`Update`** delays *everything* behind it (input feels laggy).
- Slow **`View`** hurts on every message even when nothing visible changed.
- **High message rates** (mouse motion, ticks at 1ms, firehose `Send`) multiply both costs.

Measure before optimizing: `go test -bench` on `View`, `pprof` (`net/http/pprof` on a side port, or `runtime/pprof` to a file), and `time.Since` logging into `debug.log`.

### Techniques
1. **Never block**. I/O, `time.Sleep`, heavy compute → commands (Unit 4).
2. **Window big data** (Unit 11): render only the visible rows.
3. **Cache rendered fragments** that depend on unchanged data (store the string in the model; recompute in `Update` when its inputs change):
   ```go
   case dataMsg: m.data = msg; m.tableStr = renderTable(m.data, m.w)   // once
   // View just uses m.tableStr
   ```
   Invalidate on resize as well.
4. **Build strings efficiently**: `strings.Builder`, avoid repeated `+=` in loops, avoid creating `lipgloss.NewStyle()` per row per frame (define once).
5. **Throttle noisy inputs**: `MouseModeCellMotion` not `AllMotion`; coalesce producer events (Unit 10); ignore no-op messages early (return `m, nil`).
6. **Avoid `View` doing work that belongs in `Update`**: sorting/filtering/formatting.
7. **Lip Gloss cost**: `Render` with borders/wrapping on huge strings is not free; measure `lipgloss.Width` on giant strings; render only what's visible.

### Errors as data
- **Commands return `errMsg`** (Unit 4); `Update` maps them to state (Unit 11's state enum).
- **Show errors where the user is looking**: status bar/toast for transient, full-screen for fatal; always include the recovery action ("press r to retry").
- **Timeouts everywhere** (contexts) so "loading" can't be forever.
- **Wrap with context** (`fmt.Errorf("load %s: %w", path, err)`); keep the original for `errors.Is/As`.
- **Fatal vs recoverable**: only unrecoverable startup errors should exit non-zero from `main`; runtime errors should stay in the UI.
- **Don't swallow**: log to `debug.log` even if you also show it.

### Panics and terminal restoration
Bubble Tea recovers panics in `Update`/`View`/commands by default: it restores the terminal (leaves raw mode/alt-screen), prints "Caught panic: …" with a stack trace to stderr, and `Run` returns an error wrapping `tea.ErrProgramPanic`. With `TEA_DEBUG=true` it also writes `bubbletea-panic-<unix>.log`. Consequences:
- A panic is a bug, not control flow; fix the cause (nil map, index out of range, negative `strings.Repeat`, failed type assertion).
- Don't use `WithoutCatchPanics()` unless you know why (a crash then leaves the terminal raw; recover with `reset` or `stty sane`).
- Panics in **your own goroutines** (not commands) aren't covered by the command recover wrapper; recover them yourself.

### Handing the terminal to another program
Editors, pagers, shells:
```go
c := exec.Command(editor, path)
return m, tea.ExecProcess(c, func(err error) tea.Msg { return editorFinishedMsg{err} })
```
Bubble Tea releases the terminal (leaves raw mode/alt-screen), runs the child with the real TTY, then restores and resumes your program and delivers your callback's message. For non-interactive processes use a plain `Cmd` with `exec.CommandContext` and capture output. Manual control exists (`p.ReleaseTerminal()` / `p.RestoreTerminal()`), rarely needed.

### Suspend / resume
Raw mode disables the shell's Ctrl-Z. To support it, handle the key and return `tea.Suspend`; the program is suspended and a `tea.ResumeMsg` arrives when the user foregrounds it (`fg`). Refresh stale state on resume (data may have changed).

### Printing above the UI
`tea.Println` / `tea.Printf` (commands) print a line above the managed area *in inline mode* and it persists in scrollback; in alt-screen mode they print nothing. Great for "✓ downloaded file X" logs above a progress bar.

### No TTY? Degrade gracefully
When stdin/stdout aren't terminals (pipes, CI, cron), a TUI can't work. Detect and choose a fallback (Unit 14): plain output, a `--json` mode, or run with `tea.WithoutRenderer()` (output/log behave like a normal CLI, useful for daemon-mode of the same program). Also respect `NO_COLOR`.

### Terminal-size robustness
Handle 0×0 (before first size), tiny terminals (show "terminal too small" below a minimum), and huge ones; never `strings.Repeat` with a possibly negative count (`max(0, n)`).

### Resource hygiene
Close files/connections on quit: return `tea.Quit` only after cleanup, or sequence: `tea.Sequence(saveCmd, tea.Quit)`. Cancel contexts (Unit 10). Restore anything global you changed (title is reset by the runtime; other process-level state is yours).

## Worked examples

### Straightforward: render cache
`View` was building a 2,000-line report each frame. Move formatting into `Update` on `dataMsg`/`WindowSizeMsg`, store `m.reportStr`, and `View` becomes `tea.NewView(m.reportStr)`. Cost drops from O(data) per message to O(1) per message.

### Typical: open `$EDITOR` and reload
Press `e` → `ExecProcess(vim, path)`; on `editorFinishedMsg` → return a command that re-reads the file (`loadFile`), set `loading`; on `fileLoadedMsg` show contents. See `examples/u13-exec` for the wiring (falls back to `vi` if `$EDITOR` is unset).

### Edge case: panic in a goroutine you started, and shutdown ordering
```go
go func() { p.Send(result(doWork())) }()   // doWork panics → NOT recovered by Bubble Tea's command wrapper
```
This crashes the whole process, leaving the terminal in raw mode. Wrap your goroutines with `defer recover()` → send an `errMsg`. Second edge: `return m, tea.Sequence(save, tea.Quit)` where `save` fails: it still quits (Sequence continues regardless of the returned message); if you want "quit only if saved", handle `savedMsg` in `Update` and return `tea.Quit` there, and `errMsg` shows the error instead.

## Common misconceptions
- **"Performance = FPS."** The bigger risks are blocking `Update` and heavy `View`.
- **"Panics are handled, so they're fine."** The terminal is restored but the program is dead.
- **"`fmt.Println` in a TUI prints above the UI."** Use `tea.Println` (inline only).
- **"`ExecProcess` runs concurrently."** It pauses your program while the child runs.
- **"Errors should exit the program."** Most should become UI state.
- **"Optimizing early helps."** Profile first.

## Check your understanding
1. What runs on every message? *(Update, View, renderer diff.)*
2. Where should a formatted 2,000-line report be computed? *(In `Update`, cached.)*
3. What error does `Run` return after a panic? *(One wrapping `ErrProgramPanic`.)*
4. How do you open vim from a TUI? *(`tea.ExecProcess`.)*

## Summary
- Cost per message = Update + View + render; keep them cheap, cache and window.
- Errors are data; timeouts and retry paths everywhere.
- Panics are caught and the terminal restored; fix causes; guard your own goroutines.
- `ExecProcess`, `Suspend`, `Println` integrate with the shell; degrade gracefully without a TTY.
