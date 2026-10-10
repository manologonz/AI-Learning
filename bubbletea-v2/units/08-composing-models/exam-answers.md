# Unit 8 Exam — Answers

**E1:** Add `confirmScreen`; `type confirm struct{}`; `func (c confirm) Update(msg) (confirm, tea.Cmd)` returning `func() tea.Msg { return quitConfirmedMsg{} }` on `y` and `navMsg{menuScreen}` on `n`. Menu "Quit" now returns `navMsg{confirmScreen}` instead of `tea.Quit`. Router: handle `quitConfirmedMsg` → `return m, tea.Quit`; delegate to `m.confirm` when `current == confirmScreen`; `View` renders it. Key point: the child doesn't call `tea.Quit` itself unless intended; the router owns app lifecycle.

**E2:** Causes: (1) forwarding keys to both inputs (fix: forward only to `m.inputs[m.focus]`); (2) both left `Focus()`ed (fix: `Blur()` the old one when changing focus).

**E3:** The router only forwards messages to the active child; while inactive the spinner's `TickMsg` was dropped, so its tick chain ended. Fix: forward spinner/tick messages to the child regardless of activity, or restart it on entry (`return m, child.Init()` / `spin.Tick`).

**E4:** Concrete-typed children: type-safe, no assertions, simplest for a fixed set. `tea.Model` children help when steps are pluggable; for 5 fixed different steps a step-index + concrete fields (or interface with a small `Step` interface) is fine. Best answer: concrete unless steps must be dynamically registered; justify with panics from failed assertions.

**E5:**
```go
case tea.WindowSizeMsg:
    m.w, m.h = msg.Width, msg.Height
    bodyH := m.h - 1
    m.sidebar.SetSize(24, bodyH)
    m.main.SetSize(m.w-24, bodyH)
    // footer rendered at width m.w in View
```
Guard tiny sizes: `max(0, ...)`.

**E6:** With value receivers the "parent" is a copy; the pointer refers to a stale copy or an address that no longer corresponds to the model the runtime stores, so updates through it don't affect the live model. It also hides data flow and races with commands.

**E7:**
```go
case "tab":       m.blur(); m.focus = (m.focus+1)%n; return m, m.focusCmd()
case "shift+tab": m.blur(); m.focus = (m.focus-1+n)%n; return m, m.focusCmd()
// key messages: only children[m.focus]
// all other messages (ticks, size, async): every child
```

**E8:** Costs: unreadable switch, hard to test, tangled state, merge conflicts. Benefits: no plumbing. Refactor: extract screens/widgets as child models; parent = router + layout + shared messages.

**E9:** The first `cmd` is overwritten by the second (the form's command is lost); and both children receive the key regardless of the active screen (they should be routed by `m.current`). Also `navMsg` case falls through without returning is fine in Go but no re-init on entry. Fix: switch on `m.current`, and batch commands if multiple.

**E10:** A `modal` child plus `m.modal *modal` (or `active bool`). While active: the router forwards all key/mouse messages *only* to the modal and skips the underlying screen; `View` draws the underlying view then places the modal on top via `lipgloss.Place`/layers. The modal returns its result as a message (`dialogResultMsg{id, ok}`) then the router clears it. Underlying screen still receives non-input messages (ticks, async).
