# Unit 13 Quiz — Answers

**Q1: C.** 60 by default (configurable up to 120 via `WithFPS`).

**Q2: B.** `tea.ExecProcess` pauses the program, gives the child the TTY, then resumes. A/D don't exist; C leaves the terminal in raw mode/renderer state and corrupts the display.

**Q3: B.** `Run` returns an error wrapping `ErrProgramPanic` after restoring the terminal. Not nil, not interrupt, not a re-panic.

**Q4: B.** Compute once when inputs change; `View` runs on every message. A repeats work; C/D can't react to changing data.

**Q5: B.** Documented: nothing printed when alt-screen is active.

**Q6:** A slow `View` (recomputing large strings every message); high-rate messages (`AllMotion`, fast ticks, firehose `Send`); big data rendered without windowing; expensive Lip Gloss on huge strings.

**Q7:** No: Bubble Tea's recover wrapper covers `Update`/`View` and commands, not goroutines you started; an unrecovered panic crashes the process with the terminal left raw. Wrap the goroutine body with `defer recover()` and report via `p.Send(errMsg{…})`.

**Q8:** False. `Sequence` runs commands in order regardless of message content; to make quit conditional, quit from `Update` upon `savedMsg`.

**Q9:**
```go
case tea.KeyPressMsg: if msg.String() == "ctrl+z" { return m, tea.Suspend }
case tea.ResumeMsg:    m.loading = true; return m, loadData
```

**Q10: B.** `WithoutRenderer` disables rendering so output/logging behaves like a normal CLI. `WithoutSignals` ignores OS signals; `WithFPS(0)` falls back to the default; `WithoutCatchPanics` disables recover.
