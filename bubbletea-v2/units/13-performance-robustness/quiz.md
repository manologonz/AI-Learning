# Unit 13 Quiz

**Q1 (easy).** Default maximum render FPS?
A) 10  B) 30  C) 60  D) Unlimited

**Q2 (easy).** Which function runs vim from a TUI and resumes afterwards?
A) `tea.Run`  B) `tea.ExecProcess`  C) `os/exec.Run` directly  D) `tea.Shell`

**Q3 (medium).** What does `Run` return after a recovered panic?
A) nil  B) An error wrapping `tea.ErrProgramPanic`  C) `ErrInterrupted`  D) It re-panics

**Q4 (medium).** Where is the best place to compute a heavy formatted string from data?
A) `View`  B) `Update`, cached in the model  C) `Init` only  D) `main`

**Q5 (medium).** `tea.Println` in alt-screen mode:
A) Prints above the UI  B) Prints nothing  C) Panics  D) Exits alt-screen

**Q6 (medium).** Short answer: give two ways a TUI can lag without any long `Update`.

**Q7 (hard).** Short answer: a goroutine you started panics. Is the terminal restored automatically? What do you do?

**Q8 (medium).** True/False: `tea.Sequence(save, tea.Quit)` skips the quit if `save` returns an error message.

**Q9 (hard, apply).** Show how to support Ctrl-Z and refresh data on resume.

**Q10 (medium).** Which option makes a Bubble Tea program behave like a plain CLI (no rendering)?
A) `tea.WithoutSignals()`  B) `tea.WithoutRenderer()`  C) `tea.WithFPS(0)`  D) `tea.WithoutCatchPanics()`
