# Unit 6 Quiz

**Q1 (easy).** Lip Gloss v2's import path is:
A) `github.com/charmbracelet/lipgloss`  B) `charm.land/lipgloss/v2`  C) `lipgloss.charm.sh/v2`  D) `charm.land/lipgloss`

**Q2 (easy).** `lipgloss.Color("205")` in v2 returns:
A) a string  B) `lipgloss.Color` type  C) an `image/color.Color`  D) an int

**Q3 (medium).** Order of layers from the content outward:
A) margin → border → padding  B) padding → border → margin  C) border → padding → margin  D) padding → margin → border

**Q4 (medium).** `s := lipgloss.NewStyle(); s.Bold(true); out := s.Render("x")`. Is `out` bold?
A) Yes  B) No, styles are immutable and the result was discarded  C) Only in a real TTY  D) Compile error

**Q5 (medium).** To measure a rendered, styled string's on-screen width:
A) `len(s)`  B) `utf8.RuneCountInString(s)`  C) `lipgloss.Width(s)`  D) `strings.Count(s, "")`

**Q6 (medium).** Short answer: what does `GetHorizontalFrameSize()` return?

**Q7 (hard).** A pane style has `Border(...)`, `Padding(0,1)`, `Margin(0,1)` and `Width(m.w)`, and the UI overflows the terminal by 2 columns. Why?

**Q8 (medium).** `JoinHorizontal(lipgloss.Top, a, b)` — what does `Top` do?

**Q9 (hard, apply).** Write the `View` skeleton for a 3-row layout where the body fills the remaining height between a header and footer.

**Q10 (medium).** In v2, how does a Bubble Tea app learn whether the terminal is dark or light?
