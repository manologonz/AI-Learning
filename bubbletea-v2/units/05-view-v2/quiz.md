# Unit 5 Quiz

**Q1 (easy).** In v2, `View()` returns:
A) `string`  B) `tea.View`  C) `[]string`  D) `tea.Model`

**Q2 (easy).** How do you enable alt-screen in v2?
A) `tea.WithAltScreen()`  B) `tea.EnterAltScreen` command  C) `v.AltScreen = true`  D) `p.EnterAltScreen()`

**Q3 (medium).** Which mouse mode reports hover movement with no buttons pressed?
A) `MouseModeNone`  B) `MouseModeCellMotion`  C) `MouseModeAllMotion`  D) `MouseModeHover`

**Q4 (medium).** By default, is the cursor shown in a v2 view?
A) Yes, always  B) No; `Cursor` is nil  C) Only in alt-screen  D) Only with mouse mode

**Q5 (medium).** Short answer: What does "declarative" mean for terminal features here?

**Q6 (medium).** Which v2 field delivers `FocusMsg`/`BlurMsg`?
A) `v.ReportFocus`  B) `v.Focus`  C) `tea.WithReportFocus()`  D) They always arrive

**Q7 (hard).** Short answer: your app flickers out of full-screen while `m.loading`. What's the likely bug?

**Q8 (medium).** Map each: `p.Start()`, `tea.Sequentially`, `tea.WindowSize()` → v2.

**Q9 (hard, apply).** Write a `View` that sets the window title to `"<file> (<pct>%)"` and shows the terminal progress bar only when `pct > 0`.

**Q10 (medium).** Inline (no alt-screen) vs alt-screen: which leaves output in scrollback after exit?
