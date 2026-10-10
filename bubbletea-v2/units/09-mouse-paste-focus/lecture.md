# Unit 9: Mouse, Paste, Focus & Advanced Input

> Runnable code: `examples/u09-input/main.go` — an event inspector for mouse, wheel, paste, focus and key release.

## Objectives
By the end of this unit you can:
- Enable mouse input and handle click, release, wheel, and motion messages.
- Implement hit-testing (which widget was clicked) using coordinates and layout knowledge.
- Handle bracketed paste (`PasteMsg`, `PasteStartMsg`, `PasteEndMsg`).
- React to terminal focus/blur.
- Detect and use keyboard enhancements (key release, repeat, disambiguation) safely, with fallbacks.
- Use clipboard commands (`SetClipboard`, `ReadClipboard`) and understand their limits.

## Prerequisites
Units 3, 5, 8.

## Core content

### Mouse: opt in via the View
Mouse events are **off by default**. Turn them on declaratively (Unit 5):
```go
v.MouseMode = tea.MouseModeCellMotion   // clicks, releases, wheel, drag motion
v.MouseMode = tea.MouseModeAllMotion    // + motion with no button pressed (hover)
```
Prefer `CellMotion` unless you need hover; `AllMotion` floods `Update` with messages.

Messages (v2 split them by kind; v1's single `MouseMsg` with `Action` is gone):
| Message | When |
|---|---|
| `tea.MouseClickMsg` | button pressed |
| `tea.MouseReleaseMsg` | button released |
| `tea.MouseWheelMsg` | wheel (up/down/left/right) |
| `tea.MouseMotionMsg` | movement (with/without button, per mode) |
| `tea.MouseMsg` | interface implemented by all four; `msg.Mouse()` returns the `Mouse` struct |

`Mouse` has `X, Y int` (**zero-based, origin at the top-left of the *terminal***), `Button` (`tea.MouseLeft`, `MouseRight`, `MouseMiddle`, `MouseWheelUp`, …) and `Mod` (ctrl/alt/shift held).

```go
case tea.MouseClickMsg:
    if msg.Button == tea.MouseLeft { m = m.click(msg.X, msg.Y) }
case tea.MouseWheelMsg:
    if msg.Button == tea.MouseWheelUp { m.scroll-- } else if msg.Button == tea.MouseWheelDown { m.scroll++ }
```

### Hit testing
Terminals report cell coordinates, not "which widget". You must map them. Techniques:
1. **Layout arithmetic**: you know each pane's rectangle from your layout code (Unit 6); test `x >= r.x && x < r.x+r.w && y >= r.y && y < r.y+r.h`.
2. **Row mapping**: list row = `y - topOffset + scrollOffset`.
3. **`View.OnMouse`**: an optional handler that intercepts mouse messages *based on the last rendered content* (e.g. clickable links), returning a `Cmd`. Useful when only the view knows where things are.

**Inline mode caveat:** coordinates are relative to the terminal, but in non-alt-screen mode your view sits somewhere down the screen, so `y` doesn't start at your view's top. Mouse-heavy UIs should use `AltScreen = true` so your origin is (0,0).

Never rely on the mouse alone: keyboard parity is an accessibility and SSH/multiplexer requirement, and some terminals swallow mouse events (selection needs Shift+drag when mouse mode is on).

### Bracketed paste
Terminals wrap pasted text in markers so apps can tell paste from typing. v2 delivers:
- `tea.PasteStartMsg{}` → paste begins
- `tea.PasteMsg{Content string}` → the pasted text (**one message**, not thousands of key presses)
- `tea.PasteEndMsg{}` → paste ends

```go
case tea.PasteMsg:
    m.text += msg.Content      // don't run each char through your key handlers!
```
Why it matters: pasting `"quit\n"` into a naive key handler would *execute shortcuts*. Bubbles' `textinput`/`textarea` already handle `PasteMsg`, so forward it to them. Disable with `v.DisableBracketedPasteMode = true` only if you have a reason. Sanitize pasted content (newlines in single-line inputs, control characters).

### Focus and blur
```go
v.ReportFocus = true                       // opt in (note: tmux needs focus-events on)
case tea.FocusMsg: m.focused = true        // pause? resume polling?
case tea.BlurMsg:  m.focused = false       // dim UI, pause animations, stop cursor blink
```
Uses: dim the interface, pause expensive updates, auto-refresh on return. Not every terminal/multiplexer reports this, so treat it as an enhancement, never a requirement.

### Keyboard enhancements
By default Bubble Tea requests basic **key disambiguation** from terminals that support the enhanced keyboard protocols. That's what makes `shift+enter`, `ctrl+i` vs `tab`, etc. distinguishable there. You can ask for more:
```go
v.KeyboardEnhancements.ReportEventTypes = true   // key release + repeat events
```
Then the terminal, if capable, sends a `tea.KeyboardEnhancementsMsg` telling you what it supports (`SupportsKeyDisambiguation()`, `SupportsEventTypes()`, …). Handle it:
```go
case tea.KeyboardEnhancementsMsg:
    m.canRelease = msg.SupportsEventTypes()
case tea.KeyReleaseMsg: …           // only arrives if supported & requested
case tea.KeyPressMsg:
    if msg.Key().IsRepeat { … }     // held key auto-repeat
```
**Design rule:** any feature that needs release/repeat/`shift+enter` must have a fallback path for terminals that can't deliver it, and you must not double-count (press + release) with `tea.KeyMsg` (Unit 3).

### Clipboard (OSC 52)
`tea.SetClipboard("text")` returns a `Cmd` that asks the terminal to set the system clipboard via an escape sequence; `tea.ReadClipboard` requests it (arrives as `tea.ClipboardMsg`). Support and security policy vary by terminal (many disable reads; some disable writes; it can work over SSH which is its main charm). Always show UI feedback ("Copied") only if you can tolerate it silently failing, and offer a fallback.

## Worked examples

### Straightforward: click to increment
```go
case tea.MouseClickMsg:
    if msg.Button == tea.MouseLeft { m.count++ }
```
(with `v.MouseMode = tea.MouseModeCellMotion` set in `View`).

### Typical: clickable list rows
```go
const topRow = 2                        // title + blank line above the list
case tea.MouseClickMsg:
    if msg.Button != tea.MouseLeft { break }
    idx := msg.Y - topRow + m.offset    // screen row → item index
    if idx >= 0 && idx < len(m.items) { m.cursor = idx }
case tea.MouseWheelMsg:
    switch msg.Button {
    case tea.MouseWheelUp:   m.offset = max(0, m.offset-1)
    case tea.MouseWheelDown: m.offset = min(len(m.items)-1, m.offset+1)
    }
```
The constant `topRow` is the fragile part: derive it from your layout (`lipgloss.Height(header)`) instead.

### Edge case: paste breaking a key-driven UI, and double-firing releases
Pasting `"q"` into a search field on a UI where `q` quits: if paste were delivered as key presses, the app would quit. Bracketed paste prevents this *only if* you handle `PasteMsg` in the input and don't also forward it as keys. Second edge: enabling `ReportEventTypes` while your handler is on `tea.KeyMsg` doubles every action (press+release); switch to `KeyPressMsg`. Third: on terminals without bracketed paste support, paste arrives as individual key presses; you can't fully prevent shortcut misfires there, so avoid destructive single-key shortcuts while a text field is focused.

## Common misconceptions
- **"Mouse works by default."** You must set `MouseMode`.
- **"`msg.Y` is relative to my view."** It's relative to the terminal; only equal to view-relative in alt-screen at the top-left.
- **"`tea.MouseMsg` has `Action`/`Type`."** v2 uses separate message types.
- **"Pasted text arrives as key presses."** With bracketed paste it's one `PasteMsg`.
- **"Key release events always arrive."** Only on terminals that support them and only if requested.
- **"Clipboard is guaranteed."** OSC 52 support/policy varies.

## Check your understanding
1. Which `MouseMode` includes hover? *(AllMotion.)*
2. Give the three paste-related messages. *(PasteStartMsg, PasteMsg, PasteEndMsg.)*
3. How do you learn whether release events are supported? *(`KeyboardEnhancementsMsg.SupportsEventTypes()`.)*
4. Why prefer alt-screen for mouse UIs? *(Predictable (0,0) origin.)*

## Summary
- Mouse and focus are opt-in view fields; mouse messages are split by kind with terminal-relative coordinates.
- Do hit-testing from your layout; keep keyboard parity.
- Paste arrives as a single `PasteMsg`; forward it to inputs.
- Key release/repeat need requesting and capability detection; always provide fallbacks.
