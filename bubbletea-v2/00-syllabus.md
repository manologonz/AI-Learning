# Bubble Tea v2 (Go TUI Framework) — Course Syllabus

**Learner starting level:** Go beginner, no Bubble Tea / Elm-architecture experience
**Goal:** Build and ship a real TUI application
**Depth:** Deep mastery (14 units + capstone)
**Requirements:** Go 1.25+ (auto-downloaded by the Go toolchain if needed), a real terminal
**Companion code:** every unit has a runnable, compile-checked program in `examples/` (Bubble Tea 2.0.9, Lip Gloss 2.0.6, Bubbles 2.2.1)
**Estimated total time:** ~30 hours

## Units

1. **Go Survival Kit for TUI Work** (~2h) — just the Go you need: structs, methods, interfaces, type switches, closures, modules.
   - Objectives: write a struct with methods; explain interfaces; use a type switch; write a closure; set up a module and `go get`
   - Prerequisites: none
2. **The Elm Architecture & Your First Program** (~2h) — Model, Init, Update, View; the event loop.
   - Objectives: explain the MVU loop; write a counter; explain value vs pointer receivers here; run a Program
   - Prerequisites: 1
3. **Messages & Keyboard Input (v2 key API)** (~2h) — `KeyPressMsg`, `KeyReleaseMsg`, key strings, `WindowSizeMsg`.
   - Objectives: handle keys; distinguish press/release; respond to resize; define custom messages
   - Prerequisites: 2
4. **Commands: Side Effects Without Breaking Purity** (~2.5h) — `tea.Cmd`, `Batch`, `Sequence`, `Tick`, `Quit`, async I/O.
   - Objectives: write cmds returning msgs; batch vs sequence; timers; HTTP fetch with error msg
   - Prerequisites: 3
5. **Views in v2: `tea.View`, Declarative Terminal Features** (~2h) — View struct, alt screen, cursor, mouse mode, window title set declaratively.
   - Objectives: build a `tea.View`; toggle features from View; migrate from v1 program options
   - Prerequisites: 4
6. **Styling & Layout with Lip Gloss v2** (~3h) — styles, borders, padding, colors, joins, placement, sizing.
   - Objectives: style text; compose layouts responsive to window size; handle color profiles
   - Prerequisites: 5
7. **Bubbles: Ready-Made Components** (~3h) — textinput, textarea, list, table, spinner, progress, viewport, help, key bindings.
   - Objectives: embed a component; forward msgs; use key.Binding + help
   - Prerequisites: 6
8. **Composing Models: Parent/Child & Multi-Screen Apps** (~3h) — nested models, routing, focus, state machines.
   - Objectives: design a component tree; route msgs; manage focus
   - Prerequisites: 7
9. **Mouse, Paste, Focus & Advanced Input** (~2h) — mouse click/wheel/motion messages, bracketed paste, focus events, keyboard enhancements.
   - Prerequisites: 8
10. **Concurrency & External Events** (~2.5h) — `Program.Send`, subscriptions via cmd loops, goroutines, channels, contexts, cancelation.
    - Prerequisites: 8
11. **Forms, Validation & Data-Heavy UIs** (~2.5h) — form flow, validation, tables/lists at scale, pagination.
    - Prerequisites: 8
12. **Testing & Debugging TUIs** (~2h) — testing Update as pure function, teatest, logging to file, headless runs.
    - Prerequisites: 8
13. **Performance, Errors & Robustness** (~2h) — render cost, avoiding blocking Update, panics, terminal restore, exec/suspend.
    - Prerequisites: 10
14. **Packaging & Shipping** (~1.5h) — CLI flags, config, non-TTY detection, SSH apps (Wish), releases.
    - Prerequisites: 13
15. **Capstone Build** (~4h) — design and build a complete multi-screen TUI.

## How this course works
Each unit: lecture → quiz (ungraded stakes) → unit exam → curated resources. I teach each lecture conversationally, grade your answers with per-question feedback, and log scores and weak spots in `progress.md`. A cumulative final exam and capstone close the course.

**Note:** v2 moved to vanity import paths (`charm.land/bubbletea/v2`, `charm.land/lipgloss/v2`, `charm.land/bubbles/v2`). I'll verify every API detail against current docs while writing each unit.

## What's on disk
- `units/NN-*/` — `lecture.md`, `quiz.md` + `quiz-answers.md`, `exam.md` + `exam-answers.md`, `resources.md`
- `examples/uNN-*/` — runnable, compile-checked programs (`go run ./examples/u02-counter`)
- `final-exam/` — cumulative exam (Parts A/B) + capstone task
- `progress.md` — scores and weak spots (updated as you complete assessments)

## Course status
All 15 units, the final exam, and examples are written. Suggested pace: one unit per sitting; do the quiz right after the lecture and the exam after a break.
