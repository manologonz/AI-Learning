# Final Exam — Bubble Tea v2 (cumulative)

**Time:** ~2.5 hours. Open notes/docs allowed; no running code for Part A/B (you may for Part C).
**Weighting (your goal is "build a real app"):** Part A 25%, Part B 35%, Part C (capstone task) 40%. Units 4, 7, 8, 10, 11, 12 are weighted heaviest.

## Part A — Foundations & v2 API (recall/understand/apply) — 12 questions

**A1 (Unit 1).** `func (c counter) Inc() { c.n++ }` — does `c.Inc()` change the caller's counter? Rewrite so it does.

**A2 (Unit 1).** What does a `switch m := msg.(type)` do when `msg` is `*stopMsg` but a `case stopMsg:` exists?

**A3 (Unit 2).** State the v2 `Model` interface and say what changed from v1 in `Init` and `View`.

**A4 (Unit 2).** What does `Run()` return, and why is that useful?

**A5 (Unit 3).** What are `String()` and `Key().Text` for: space, shift+a, ctrl+c, up-arrow?

**A6 (Unit 3).** Why can `case tea.KeyMsg:` double-fire?

**A7 (Unit 4).** `Batch` vs `Sequence`: choose for (i) save then quit, (ii) fetch two endpoints in parallel.

**A8 (Unit 4).** Why does a clock stop after the first tick, and how do you fix it?

**A9 (Unit 5).** Convert to v2: `NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())`.

**A10 (Unit 6).** In Lip Gloss v2, a style has `Border`, `Padding(0,2)`, `Margin(0,1)`, `Width(50)`. What is the rendered width and the content width?

**A11 (Unit 7).** What are the three steps of embed–forward–batch, and what breaks if you skip the returned `Cmd`?

**A12 (Unit 9).** Give the three paste messages and one reason bracketed paste matters.

## Part B — Design, analysis & debugging (analyze/evaluate) — 8 questions

**B1 (Unit 8).** A form and a list are both on screen. Typing in the form also moves the list cursor. Diagnose and fix (write the routing code).

**B2 (Unit 4/10).** Search-as-you-type: replies arrive out of order and an old reply overwrites newer results. Give a complete fix including cancellation.

**B3 (Unit 10).** A goroutine reads `m.count` for a debug endpoint. Explain the bug and two safe designs.

**B4 (Unit 11).** After a refresh, the highlighted row jumps. Root cause and fix. Also explain why `loading`/`err` booleans are a hazard and give the alternative.

**B5 (Unit 12).** Write (a) an `Update` test and (b) a teatest flow for "press `a` opens form; type `hi`; press `enter`; list shows `hi`". State what you fake and how you avoid hangs.

**B6 (Unit 13).** The UI lags on a 50k-row view although `Update` is fast. Explain three plausible causes and fixes; also explain what a panic in your own goroutine does.

**B7 (Unit 14).** Design the layout and entry points for a tool that supports TUI, `--json`, and SSH; explain per-session state and one security consideration.

**B8 (Units 5, 6, 8).** "Terminal too small" and full-screen toggling: describe how `View` should handle sizes and `AltScreen` so state isn't lost and there's no flicker.

## Part C — Capstone synthesis (create)

**C1.** Build (or extend your Taskmate from Unit 15) so that it demonstrates *all* of the following, and submit code + a ≤1-page design note:
1. Explicit data-state enum with async load/save via commands, timeouts, and error+retry.
2. Router with at least three screens and a modal; child→parent communication via messages.
3. A `list` with filter, selection by ID, and a form with validation (including "untouched required field" handling).
4. Responsive layout using v2 width rules, a tiny-terminal guard, and light/dark styles.
5. Sequence-ID based stale-message protection (toast or search).
6. Tests: table-driven `Update` tests, at least one `View` assertion, one teatest flow with timeouts; `go test -race` clean.
7. `--version`, exit codes, non-TTY fallback.

*(Grading is by the rubric in Unit 15 plus the checklist above.)*
