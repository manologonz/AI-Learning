# Unit 6: Styling & Layout with Lip Gloss v2

> Runnable code: `examples/u06-layout/main.go` — resize your terminal and watch the panes reflow.

## Objectives
By the end of this unit you can:
- Create and chain `lipgloss.Style` values (color, bold, padding, margin, border, width/height).
- Compose layouts with `JoinHorizontal`, `JoinVertical`, and `Place`.
- Size panes from `tea.WindowSizeMsg` without off-by-border bugs (frame sizes).
- Use the v2 color model (`lipgloss.Color(...)` returns a `color.Color`) and adapt light/dark.
- Explain that styles are immutable values and why "render" produces a string.

## Prerequisites
Unit 5.

## Core content

### Import
```go
import "charm.land/lipgloss/v2"     // v2 (v1 was github.com/charmbracelet/lipgloss)
```

### Styles are immutable values
```go
var title = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#FAFAFA")).
    Background(lipgloss.Color("#7D56F4")).
    Padding(0, 1)

s := title.Render("Hello")     // returns a string containing ANSI codes
```
Each builder method returns a **new copy**, so you can safely derive: `warn := title.Background(lipgloss.Color("196"))`. `Render` returns a string, so styles compose with everything from Units 2–5: the string goes into `tea.NewView`.

**v2 gotcha:** `lipgloss.Color("205")` is now a *function* returning `image/color.Color`, not a string type. So `var c lipgloss.Color = "205"` (v1) no longer compiles; write `var c color.Color = lipgloss.Color("205")` or just pass it inline.

Color forms: ANSI 0–15, ANSI256 (`"205"`), or hex (`"#7D56F4"`). Bubble Tea/Lip Gloss downsample automatically to what the terminal supports (Bubble Tea handles this for you; standalone programs use `lipgloss.Println`).

### Spacing and borders
```go
box := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).      // Normal, Thick, Double, Hidden…
    BorderForeground(lipgloss.Color("63")).
    Padding(1, 2).                         // CSS-like: (v,h) or (t,r,b,l)
    Margin(1)
```
Order (from inside out): **content → padding → border → margin.** Padding is inside the border; margin outside.

### Width and height (the #1 source of layout bugs)
In Lip Gloss v2, `Width(n)` / `Height(n)` set the size of the block **including padding and border**, but **not margin** (margin is added outside). Verified on Lip Gloss 2.0.6:

```go
s := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1, 2).Width(20)
lipgloss.Width(s.Render("hi"))    // 20  (border and padding are INSIDE the 20)
// content area = 20 - 2 (border) - 4 (padding) = 14 columns
```
With a margin, the total grows: `.Margin(1, 2).Width(20)` renders 24 columns wide (20 + 2×2 margin).

So:
- To fill the terminal with a bordered pane: `Width(m.w)` is exactly right (no border subtraction).
- To know the **content** width (for wrapping inner text or sizing a `viewport`): `m.w - style.GetHorizontalFrameSize()`; frame = border + padding + margin.
- If the style has a **margin**, subtract `GetHorizontalMargins()` from the `Width` you pass, or the block overflows.

Never hard-code "minus 2 for the border"; ask the style, so a design change doesn't break your math. `Height(n)` **pads** short content up to `n` but does **not clip** taller content (verified); use `MaxHeight`/`MaxWidth` to truncate, or a `viewport` (Unit 7). Text longer than the content width wraps. Measure rendered strings with `lipgloss.Width(s)` / `lipgloss.Height(s)`, never `len(s)` (ANSI codes and wide runes break it).

> **Migration warning:** v1 `Width` excluded the border (border was extra). Layout code ported from v1 that subtracted border sizes may now be too narrow; code that didn't may have been relying on v1 behavior. Re-derive from the rule above and test.

### Composing blocks
```go
row  := lipgloss.JoinHorizontal(lipgloss.Top, left, right)        // align tops
page := lipgloss.JoinVertical(lipgloss.Left, header, row, footer) // align left edges
full := lipgloss.Place(m.w, m.h, lipgloss.Center, lipgloss.Center, page)   // center in a box
```
Position args: `lipgloss.Top/Center/Bottom` (0…1 floats) for horizontal joins; `Left/Center/Right` for vertical joins/Place.

### Making layouts responsive
Store `w, h` from `WindowSizeMsg` (Unit 3), compute pane sizes **in `View`** (or in `Update` on resize) from those numbers. Guard the zero-size first frame.
```go
half := m.w / 2
left  := pane.Width(half).Render(...)
right := pane.Width(m.w - half).Render(...)     // give the remainder to one side: no lost column on odd widths
```

### Light vs. dark terminals (v2)
v1 auto-detected the background globally. v2 removes hidden global state: you either ask the terminal in Bubble Tea or use helpers.
```go
// Init: ask the terminal for its background color
func (m model) Init() tea.Cmd { return tea.RequestBackgroundColor }
// Update:
case tea.BackgroundColorMsg:
    m.isDark = msg.IsDark()
// then choose colors: 
lightDark := lipgloss.LightDark(m.isDark)
fg := lightDark(lipgloss.Color("#333333"), lipgloss.Color("#DDDDDD"))
```
Bubbles components take an `isDark` argument for their default styles (Unit 7). The `compat` package exists for quick v1-style migration but uses global state.

### Other pieces
- `lipgloss/table`, `lipgloss/list`, `lipgloss/tree` subpackages render static tables/lists/trees to strings (not interactive; Bubbles `table`/`list` are the interactive ones).
- `Style.Inline(true)` ignores newlines/margins for single-line pieces; `Transform(fn)` post-processes text; `Underline`, `Italic`, `Faint`, etc.

## Worked examples

### Straightforward: a badge
```go
badge := lipgloss.NewStyle().Bold(true).Padding(0, 1).
    Foreground(lipgloss.Color("#000")).Background(lipgloss.Color("#FFD700"))
fmt.Println(badge.Render("NEW"))   // (in a non-Bubble-Tea program use lipgloss.Println)
```

### Typical: header / body / footer that fills the terminal
```go
func (m model) View() tea.View {
    if m.w == 0 { return tea.NewView("") }
    header := headerStyle.Width(m.w).Render("My App")
    footer := footerStyle.Width(m.w).Render("q quit")
    bodyH := m.h - lipgloss.Height(header) - lipgloss.Height(footer)
    body := bodyStyle.Width(m.w).Height(bodyH).Render(m.content)
    v := tea.NewView(lipgloss.JoinVertical(lipgloss.Left, header, body, footer))
    v.AltScreen = true
    return v
}
```
Key trick: **measure the rendered header/footer** instead of assuming they're 1 line; the body gets what's left.

### Edge case: margin overflows, `Height` doesn't clip
```go
p := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Margin(0, 1).Width(m.w).Height(m.h)
```
The border is inside `Width`, but the **margin is outside**: rendered width is `m.w + 2`, overflowing the terminal by 2 columns, which the terminal wraps and the layout shreds. Fix: `Width(m.w - p.GetHorizontalMargins())` and `Height(m.h - p.GetVerticalMargins())`. Also: if `m.content` has more lines than the pane, `Height` will **not** cut it; the block just grows taller than the screen. Clip with `MaxHeight`, or render the body through a `viewport` sized to the pane's *content* height (`m.h - p.GetVerticalFrameSize()`).

## Common misconceptions
- **"`Color` is a type."** In v2 it's a function returning `color.Color`.
- **"`Width(m.w)` leaves room for the border."** In v2 the border is *inside* `Width`; the **margin** is what sticks out. (v1 behaved differently.)
- **"Styles mutate."** They're immutable; `s.Bold(true)` alone discards the result.
- **"`len(s)` measures width."** Use `lipgloss.Width`; ANSI codes and wide characters distort `len`.
- **"Padding is outside the border."** It's inside; margin is outside.
- **"Dark/light detection is automatic."** In v2 you ask (`RequestBackgroundColor`) or pass `isDark`.

## Check your understanding
1. What does `title.Bold(true)` do if you ignore its return value? *(Nothing.)*
2. `Width(40)` with a normal border, 1-column padding and no margin: total rendered width and content width? *(40 total; content = 40 − 2 border − 2 padding = 36.)*
3. Which function gives border+padding+margin horizontally? *(`GetHorizontalFrameSize`)*
4. Why prefer `lipgloss.Height(header)` over the constant 1? *(Header text/style may wrap or change.)*

## Summary
- Styles are immutable builders that `Render` to strings; v2 colors are `color.Color` values.
- Padding inside border, margin outside; in v2 `Width`/`Height` include padding **and border**, but not margin; `Height` pads but doesn't clip.
- Derive sizes with `Get*FrameSize` and `lipgloss.Width/Height`; join blocks with `Join*` and `Place`.
- Resize layouts from `WindowSizeMsg`; v2 asks the terminal for dark/light (`BackgroundColorMsg`).
