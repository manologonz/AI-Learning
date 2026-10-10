# Unit 5 Quiz — Answers

**Q1: B.** A is v1. C/D are not the signature.

**Q2: C.** A/B/D are the removed v1 APIs.

**Q3: C.** `AllMotion` includes motion without pressed buttons. B reports clicks/drags/wheel only. D doesn't exist.

**Q4: B.** `Cursor` nil = hidden. You opt in with `tea.NewCursor`.

**Q5:** You describe the terminal state you want as data returned from `View`, derived from the model; the runtime figures out and applies the escape sequences, and reverts things you stop requesting. No enable/disable commands to sequence.

**Q6: A.** C was removed with v1's options; D is false; B doesn't exist.

**Q7:** One `View` branch (the loading one) returns a view without `AltScreen = true`, so the desired state changes. Set fields once for all branches.

**Q8:** `p.Start()` → `p.Run()`; `tea.Sequentially` → `tea.Sequence`; `tea.WindowSize()` → `tea.RequestWindowSize` (returns `Msg`, so it's directly a `Cmd`).

**Q9:**
```go
v := tea.NewView(body)
v.WindowTitle = fmt.Sprintf("%s (%d%%)", m.file, m.pct)
if m.pct > 0 { v.ProgressBar = tea.NewProgressBar(tea.ProgressBarDefault, m.pct) }
return v
```

**Q10:** Inline. Alt-screen restores the previous screen content on exit, so the app's output disappears.
