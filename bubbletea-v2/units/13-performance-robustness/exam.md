# Unit 13 Exam

**E1 (Analyze).** A log viewer renders 50k lines in `View`. Typing lags. Propose a redesign with specifics (data structure, caching, windowing).

**E2 (Apply).** Wrap a background worker goroutine so a panic becomes an `errMsg` in the UI.

**E3 (Evaluate).** Choose error UX for: (a) config file missing at startup; (b) HTTP timeout mid-session; (c) invalid user input; (d) disk full while saving.

**E4 (Apply).** Implement `open in $EDITOR then reload`: messages, commands, and state (including error handling if the editor exits non-zero).

**E5 (Analyze).** `View` calls `lipgloss.NewStyle()...` 200 times per frame. Is that a problem? Measure how you'd decide.

**E6 (Evaluate).** When would you use `tea.WithoutCatchPanics()`?

**E7 (Apply).** Add a "terminal too small" screen when width < 40 or height < 10, keeping state intact so the UI resumes after resize.

**E8 (Analyze).** After returning from `ExecProcess`, the screen shows old content in the wrong place. What could cause this and what's the right response?

**E9 (Evaluate).** Design shutdown for a program with a database connection, a running download, and unsaved edits.

**E10 (Create).** Write a checklist ("robustness review") of 8 items to run against any Bubble Tea app before release.
