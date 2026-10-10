# Unit 6 Quiz — Answers

**Q1: B.** v2 uses the vanity path with `/v2`. A is v1. C/D are invented/incomplete.

**Q2: C.** `color.Color` from the standard library; A/B are v1 behavior.

**Q3: B.** Padding is inside the border; margin is outside.

**Q4: B.** Builder methods return new styles; `s` was never reassigned. (Watch for this in code review.)

**Q5: C.** It understands ANSI sequences and cell widths. A counts bytes incl. escape codes; B counts runes incl. escapes; D counts positions.

**Q6:** The total horizontal space taken by border, padding and margin: subtract from available width to get the content width.

**Q7:** In v2 `Width` already includes border and padding, but the **margin** (1 column each side) is added outside, so the block is `m.w + 2` wide. Subtract `GetHorizontalMargins()` from the width you pass.

**Q8:** Aligns the blocks' top edges when heights differ (vs Center/Bottom).

**Q9:**
```go
header := hStyle.Width(m.w).Render("Title")
footer := fStyle.Width(m.w).Render("help")
bodyH := m.h - lipgloss.Height(header) - lipgloss.Height(footer)
body := bStyle.Width(m.w).Height(bodyH).Render(m.content)
out := lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
```
(guard `m.w == 0` first).

**Q10:** Request it with `tea.RequestBackgroundColor` (e.g. from `Init`) and handle `tea.BackgroundColorMsg` (`msg.IsDark()`), then choose colors with `lipgloss.LightDark(isDark)` or pass `isDark` into Bubbles style helpers.
