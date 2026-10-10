# Unit 15: Capstone Build — "Taskmate", a multi-screen TUI

This unit is a guided build, not a lecture. You apply every previous unit to ship a small but real application. A compiling starter is in `examples/u15-capstone/`.

## Objectives
By the end you can:
- Design an app from requirements: state machine, message vocabulary, component tree.
- Implement it with a router, child models, Bubbles components, Lip Gloss layout, commands, and persistence.
- Test it at all four levels (Unit 12) and package it (Unit 14).

## Prerequisites
Units 1–14.

## The product
**Taskmate** — a keyboard-and-mouse task manager.

### Requirements (MVP)
1. **List screen**: tasks in a `list` (title, status, due date), filter with `/`, toggle done with `space`, delete with `d` (with confirm), add with `a`.
2. **Add/Edit form**: title (required, ≤ 60 chars), notes (multi-line `textarea`), due date (`YYYY-MM-DD`, optional, validated).
3. **Detail pane**: on wide terminals (≥ 100 cols), show the selected task's notes in a `viewport` beside the list; on narrow terminals, `enter` opens the detail as a full screen.
4. **Persistence**: tasks saved as JSON in `os.UserConfigDir()/taskmate/tasks.json` using **atomic writes**; loading happens asynchronously at startup with loading / empty / error / ready states.
5. **Status bar**: shows counts (open/done), transient messages ("Saved", errors with retry).
6. **Help**: `?` toggles full help generated from `key.Binding`s.
7. **Theming**: light/dark aware via `tea.RequestBackgroundColor`.
8. **Mouse**: click selects a row; wheel scrolls; keyboard parity for everything.
9. **Resizing & tiny terminals**: handle any size, show "terminal too small" under 40×10.
10. **Exit**: `q`/`ctrl+c`; if there are unsaved changes, confirm via a modal; on quit save first (`Sequence`).

### Stretch goals
- `taskmate list --json` (non-TUI path); non-TTY fallback.
- Undo for delete (`u`), using an undo stack of previous `[]Task` snapshots (careful with slice aliasing, Unit 2).
- Serve over SSH with Wish (per-session store!).
- Recurring tasks / due-date colors (overdue red, today yellow).

## Architecture guide (a suggested decomposition)
```
internal/core      Task, Store interface, Filter/Sort, validation (no Bubble Tea)
internal/storage   JSON file store, atomic write
internal/ui
   app.go          root model: router, window size, global keys, modal, status bar
   list.go         list screen (bubbles/list + optional viewport detail)
   form.go         add/edit form (textinput, textarea, validation)
   confirm.go      modal
   keys.go         key.Binding maps + help
   styles.go       lipgloss styles from isDark
cmd/taskmate       main: flags, TTY detection, wiring
```

### Message vocabulary (draft, refine it)
`tasksLoadedMsg{tasks}`, `loadErrMsg{err}`, `savedMsg{}`, `saveErrMsg{err}`, `openFormMsg{task *core.Task}`, `formSubmittedMsg{task}`, `formCancelledMsg{}`, `confirmMsg{id, ok}`, `statusMsg{text, isErr}`, `clearStatusMsg{id}` (with sequence ID to expire toasts correctly), `tea.WindowSizeMsg`, `tea.BackgroundColorMsg`.

### State machine
`screen`: `list | form | detail`; overlay: `modal *confirm`; data state: `loading | error | empty | ready` (Unit 11).

## Build plan (checkpoints)
1. **Core + tests**: `core.Task`, validation, filter/sort; table-driven tests. *(Units 1, 12)*
2. **Skeleton app**: root model, window size, tiny-terminal guard, quit. *(Units 2, 3, 6)*
3. **Async load/save**: Store interface, commands, state enum, error/retry. *(Units 4, 11, 13)*
4. **List screen**: `list` with custom items, filter, toggle, add/delete. *(Unit 7)*
5. **Form screen**: fields, focus, validation, submit/cancel messages. *(Units 7, 8, 11)*
6. **Detail + responsive layout**: viewport, wide/narrow policies. *(Units 6, 8)*
7. **Modal + status bar + help**: confirm delete/quit, toasts with expiry IDs. *(Units 8, 4)*
8. **Mouse + theme**: click/wheel, `isDark`. *(Units 6, 9)*
9. **Tests**: Update tests per screen, one teatest end-to-end flow, `-race`. *(Units 10, 12)*
10. **Ship**: version flag, config path, non-TTY fallback, README + keybinding table, cross-compile. *(Unit 14)*

## Evaluation rubric (100 points)
| Area | Points | Excellent looks like |
|---|---|---|
| Architecture & separation | 20 | core has no Bubble Tea imports; clear router/child/message design; one source of truth |
| Correctness & UX | 20 | all MVP requirements work; keyboard parity; sensible empty/error/loading states |
| Async & robustness | 15 | no blocking in Update; timeouts; stale-reply guards; panics/goroutines handled; atomic writes |
| Layout & styling | 10 | responsive, correct width math (v2 rules), light/dark, tiny-terminal guard |
| Testing | 15 | pure-logic + Update + View + one teatest flow; fakes for I/O; `-race` clean |
| Packaging & docs | 10 | flags, exit codes, non-TTY behavior, README, keybinding help |
| Code quality | 10 | small files, idiomatic Go, comments where non-obvious |

## Worked example: designing the toast expiry (a tricky part)
Problem: "Saved" should vanish after 3s, but a second message shown 1s later must not be cleared by the first message's timer.
```go
type statusMsg struct{ text string }
type clearStatusMsg struct{ id int }

case statusMsg:
    m.statusID++
    m.status = msg.text
    id := m.statusID
    return m, tea.Tick(3*time.Second, func(time.Time) tea.Msg { return clearStatusMsg{id} })
case clearStatusMsg:
    if msg.id == m.statusID { m.status = "" }       // ignore stale timers (Unit 4)
```
This single pattern (sequence IDs) recurs in search, debounce, and toasts.

## Common capstone mistakes
- Putting file I/O in `Update`.
- Storing selected *index* instead of ID.
- Letting the list capture `q`/`esc` in filter mode while global shortcuts fire.
- Using v1 idioms (`KeyMsg`, `" "`, `tea.WithAltScreen`).
- Forgetting to re-init a form on entry.
- Hard-coded border/padding math from v1.

## How to work
Build in the checkpoints above; after each, run the tests and the app. When stuck, re-read the relevant unit's "Common misconceptions". Ask me to review any checkpoint; I'll respond as a code reviewer using the rubric.
