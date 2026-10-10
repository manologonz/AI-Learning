# Unit 6 Exam — Answers

**E1:**
```go
bodyH := m.h - 1                                   // minus footer
left  := pane.Width(30).Height(bodyH).Render(m.left)   // v2: border is INSIDE Width/Height
right := pane.Width(m.w - lipgloss.Width(left)).Height(bodyH).Render(m.right)
row := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
out := lipgloss.JoinVertical(lipgloss.Left, row, footer)
```
Key: v2 `Width`/`Height` include the border, so no border subtraction is needed for outer size; use the *measured* left width so the right pane gets exactly the remainder. (Content widths inside the panes are `Width − GetHorizontalFrameSize()`.) Answers that subtract borders are v1-style and leave the layout 2 cells short.

**E2:** `81/2 = 40` twice = 80: one column is lost (a visible gap at the right edge). Give the remainder to one side: `right.Width(m.w - half)`.

**E3:** (a) 20: in v2 both border and padding are inside `Width`. (b) content = 20 − 2 (border) − 4 (padding 2+2) = 14. Adding `Margin(0,3)` adds 3 on each side outside: total 26. (Verified against Lip Gloss 2.0.6; v1 would have rendered 22 for (a).)

**E4:** `Init: tea.RequestBackgroundColor`; `case tea.BackgroundColorMsg: m.isDark = msg.IsDark()`; then `ld := lipgloss.LightDark(m.isDark); fg := ld(lipgloss.Color("#222"), lipgloss.Color("#EEE"))`.

**E5:** `len` counts ANSI escape bytes and multibyte runes, so padding is wrong. Use `pad := w - lipgloss.Width(label)` (clamped ≥ 0), or render with `lipgloss.NewStyle().Width(w).Align(lipgloss.Right)`.

**E6:** Styles are immutable and cheap, but rebuilding them every frame allocates for nothing; define once. (Exception: styles that depend on dynamic width can be *derived* in View from the stored base via `.Width()`.)

**E7:**
```go
l, r := "main.go", "Ln 4, Col 2"
gap := max(1, m.w-lipgloss.Width(l)-lipgloss.Width(r))
bar := bar.Render(l + strings.Repeat(" ", gap) + r)
```

**E8:** `lipgloss/table` (renders a string; no interactivity, works with plain stdout via `lipgloss.Println`). The Bubbles `table` is a Bubble Tea component needing a running program, so overkill for static output.

**E9:** Bubble Tea's renderer downsamples colors and controls the writer; a plain `fmt.Println` of styled text in a normal program doesn't, so `lipgloss.Println` (and friends) writes through a color-profile-aware writer so output degrades correctly (e.g. when piped or on 16-color terminals).

**E10:** Fields: `w,h`, `sidebarOpen bool`. `View`: `if m.w < 20 { return tea.NewView("terminal too small") }`; `showSidebar := m.w >= 100 || m.sidebarOpen`; if showSidebar and `m.w >= 100` use `JoinHorizontal(sidebar(30), main(m.w-30))`; if narrow and toggled, show the sidebar *instead of* main (or overlay). Key handling: `tab` toggles `sidebarOpen`; on resize to ≥100, ignore the toggle. Good answers handle every width range explicitly.
