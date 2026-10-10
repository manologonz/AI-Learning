# Unit 5: Views in v2 — `tea.View` and Declarative Terminal Features

> Runnable code: `examples/u05-view/main.go` — toggle alt screen, mouse, focus reports, progress bar live.

## Objectives
By the end of this unit you can:
- Build a `tea.View` with `tea.NewView`, `SetContent`, and field assignments.
- Enable/disable alt-screen, mouse mode, focus reporting, bracketed paste, window title, cursor, colors, and the terminal progress bar *from `View`*.
- Explain why v2 made these features declarative and what v1 imperative APIs they replace.
- Show a real cursor (for text input) with `tea.Cursor`.
- Migrate a v1 program's options and commands to v2 fields.

## Prerequisites
Units 2–4.

## Core content

### The big idea
In v1 you toggled terminal features by *commands and startup options*: `tea.WithAltScreen()`, `tea.EnterAltScreen`, `tea.EnableMouseCellMotion`, … These scattered state across `NewProgram`, `Init` and `Update`, and it was easy to get "I entered alt screen but forgot to leave" bugs.

In v2, `View()` returns a **struct** describing *everything*: the content **and** the terminal modes you want *right now*. Bubble Tea diffs the desired state against the current one and applies the changes. It is the same MVU idea applied to terminal modes: *derive them from the model*.

```go
func (m model) View() tea.View {
    v := tea.NewView("Hello!")
    v.AltScreen = true
    v.MouseMode = tea.MouseModeCellMotion
    v.WindowTitle = "my app"
    return v
}
```

### The fields
| Field | Purpose |
|---|---|
| `Content` | the rendered string (set via `NewView(s)` or `v.SetContent(s)`) |
| `AltScreen bool` | full-screen mode; terminal restores your scrollback on exit |
| `MouseMode` | `tea.MouseModeNone` (default) / `MouseModeCellMotion` (clicks, drags, wheel) / `MouseModeAllMotion` (also hover motion) |
| `ReportFocus bool` | deliver `tea.FocusMsg` / `tea.BlurMsg` |
| `DisableBracketedPasteMode bool` | (paste events are ON by default) |
| `WindowTitle string` | terminal window/tab title |
| `Cursor *tea.Cursor` | position/shape/color/blink; `nil` = hidden |
| `ForegroundColor`, `BackgroundColor` | set terminal default colors |
| `ProgressBar *tea.ProgressBar` | native terminal progress indicator (taskbar/tab), where supported |
| `KeyboardEnhancements` | request key release/repeat etc. (Unit 9) |
| `OnMouse` | intercept mouse events based on what's on screen (Unit 9) |

Because these are *fields on a returned value*, they're **functions of the model**:
```go
v.AltScreen = m.fullscreen            // toggled by a key press in Update
v.WindowTitle = fmt.Sprintf("%s — %d%%", m.file, m.pct)
if m.pct > 0 { v.ProgressBar = tea.NewProgressBar(tea.ProgressBarDefault, m.pct) }
```
No command, no toggle: `Update` flips a bool; `View` reflects it.

### The cursor
By default no cursor is shown (`Cursor == nil`). For text input, either draw a fake cursor glyph yourself or ask for the real one:
```go
v.Cursor = tea.NewCursor(x, y)     // x,y relative to the top-left of your content
// customize:
v.Cursor.Shape = tea.CursorBar
v.Cursor.Blink = false
```
Bubbles' `textinput`/`textarea` compute cursor positions for you (Unit 7). Remember: *you* must compute where the cursor belongs; the row/column are in terminal cells within your view.

### v1 → v2 mapping (memorize the pattern)
| v1 | v2 |
|---|---|
| `tea.NewProgram(m, tea.WithAltScreen())` | `v.AltScreen = true` |
| `tea.WithMouseCellMotion()` / `EnableMouseCellMotion` | `v.MouseMode = tea.MouseModeCellMotion` |
| `tea.WithReportFocus()` | `v.ReportFocus = true` |
| `tea.SetWindowTitle("x")` | `v.WindowTitle = "x"` |
| `tea.HideCursor` / `ShowCursor` | `v.Cursor = nil` / `tea.NewCursor(...)` |
| `p.Start()` / `StartReturningModel()` | `p.Run()` |
| `tea.WindowSize()` | `tea.RequestWindowSize` |
| `tea.Sequentially` | `tea.Sequence` |
| `View() string` | `View() tea.View` |

Program options that **remain**: `WithContext`, `WithInput`, `WithOutput`, `WithEnvironment`, `WithFilter`, `WithFPS`, `WithoutSignalHandler`, `WithoutCatchPanics`, `WithoutRenderer`, plus new `WithColorProfile` and `WithWindowSize` (both great for tests).

### Inline vs. alt-screen
Without `AltScreen`, your UI renders *inline* in the normal terminal buffer (good for prompts, progress, short pickers). The output stays in scrollback after exit. With `AltScreen`, you get a clean full-screen canvas that vanishes on exit (good for dashboards/editors). Pick based on whether the output should persist.

### Rendering is cell-based and diffed
v2's renderer (the "cursed renderer") keeps a screen buffer and only writes what changed, so `View` may be called freely; making `View` **cheap and side-effect free** is your job (Unit 13).

## Worked examples

### Straightforward: a titled full-screen app
```go
func (m model) View() tea.View {
    v := tea.NewView(m.render())
    v.AltScreen = true
    v.WindowTitle = "Inbox"
    return v
}
```

### Typical: a feature toggled at runtime
```go
case tea.KeyPressMsg:
    if msg.String() == "m" { m.mouse = !m.mouse }
// View:
if m.mouse { v.MouseMode = tea.MouseModeCellMotion }
```
Pressing `m` flips the bool; next `View` requests mouse mode; Bubble Tea sends the terminal the enable sequence. Press again → mode reverts to `None` automatically. In v1 you'd have needed a command each way.

### Edge case: setting a mode in only one code path
```go
func (m model) View() tea.View {
    if m.loading {
        return tea.NewView("Loading…")            // AltScreen = false here!
    }
    v := tea.NewView(m.body())
    v.AltScreen = true
    return v
}
```
Every `View` call must describe the **complete** desired state. The `loading` branch returns a view with `AltScreen == false`, so the app flickers out of full-screen while loading and back in afterward. Build one `tea.View`, set the fields once, and vary only the content.

## Common misconceptions
- **"`View()` still returns a string."** In v2 it returns `tea.View`.
- **"I toggle alt-screen with a command."** No. Those commands were removed; set the field.
- **"`tea.WithAltScreen()` still works."** Removed in v2.
- **"Fields I don't set stay as they were."** Each `View` is a full description; omitted = default (off).
- **"Mouse mode is on by default."** Default is `MouseModeNone`.
- **"Setting `Cursor` moves it relative to the terminal."** Position is relative to your view's top-left.

## Check your understanding
1. How do you hide the cursor in v2? *(`v.Cursor = nil`, the default.)*
2. Your app leaves full-screen only during loading. What's the likely bug? *(Different `View` branch missing `AltScreen = true`.)*
3. What replaced `tea.WithMouseCellMotion()`? *(`v.MouseMode = tea.MouseModeCellMotion`.)*
4. Difference between `MouseModeCellMotion` and `AllMotion`? *(AllMotion also reports mouse movement with no button pressed.)*

## Summary
- `View()` returns a `tea.View` describing content *and* terminal modes.
- Most v1 startup options/commands became `View` fields; toggle them via model state.
- Set fields consistently in every code path; omitted fields revert to defaults.
- Inline vs alt-screen is a product decision about scrollback.
