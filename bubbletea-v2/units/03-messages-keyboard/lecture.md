# Unit 3: Messages & Keyboard Input (the v2 key API)

> Runnable code: `examples/u03-keys/main.go` — run it and press keys to *see* what each message contains.

## Objectives
By the end of this unit you can:
- Explain what a `tea.Msg` is and define your own message types.
- Handle key presses with `tea.KeyPressMsg`, using `String()` and the `Key` struct fields.
- Distinguish `KeyPressMsg`, `KeyReleaseMsg` and the `KeyMsg` interface.
- Match modifier combos (`ctrl+c`, `shift+tab`, `alt+x`) and special keys (`enter`, `esc`, `space`).
- Respond to `tea.WindowSizeMsg` and store the terminal size in the model.

## Prerequisites
Unit 2.

## Core content

### Everything is a message
`tea.Msg` is (in v2) an alias of an event interface from the underlying `ultraviolet` library, and in practice **any Go value can be delivered to `Update`**. Built-in messages include key presses, mouse events, resizes, focus/blur, pastes; you add your own (a timer fired, an HTTP reply arrived). Your `Update` type-switches on them (Unit 1).

Defining a custom message is just defining a type:
```go
type clearMsg struct{}
type statusMsg struct{ code int }
type errMsg struct{ err error }
```
Use distinct *types* for distinct meanings, not strings, so the type switch stays exhaustive and typo-proof.

### Key presses in v2
v1 had one `tea.KeyMsg` struct. v2 splits it:

| Type | Meaning |
|---|---|
| `tea.KeyPressMsg` | a key went down (what you want 95% of the time) |
| `tea.KeyReleaseMsg` | a key came up (only if the terminal supports it and you asked, Unit 9) |
| `tea.KeyMsg` | an **interface** both satisfy; use it to catch both |

```go
case tea.KeyPressMsg:
    switch msg.String() {
    case "q", "ctrl+c": return m, tea.Quit
    case "enter":       ...
    case "esc":         ...
    case "space":       ...   // NOT " " in v2
    case "shift+tab":   ...
    case "alt+x":       ...
    }
```

`String()` returns a human-readable name: printable keys give the character (`"a"`, `"A"`, `"?"`), specials give names (`"enter"`, `"tab"`, `"backspace"`, `"up"`, `"pgdown"`, `"f5"`), and modifiers are prefixed in the order `ctrl+alt+shift+…`.

### The `Key` struct (when strings aren't enough)
`msg.Key()` returns a `tea.Key`:
- `Code rune`: the key: `'a'`, or a constant like `tea.KeyEnter`, `tea.KeyUp`.
- `Text string`: the text the key **types**, e.g. `"a"`, or `""` for keys that don't type anything (arrows, enter, ctrl combos).
- `Mod`: modifiers (`tea.ModCtrl`, `tea.ModShift`, `tea.ModAlt` …).
- `ShiftedCode`, `BaseCode`, `IsRepeat`: newer info available on terminals with keyboard enhancements (Unit 9).

Use `Text` when you're building a text editor: append `msg.Key().Text` to a buffer (it's empty for non-typing keys, so arrows won't insert junk). Use `String()` for shortcuts.

Field-based matching also works:
```go
if k := msg.Key(); k.Code == 'c' && k.Mod == tea.ModCtrl { ... }
```
`Keystroke()` is like `String()` but always includes modifier info (e.g. `"shift+a"` vs `"A"`).

### Terminal reality check
Terminals historically can't distinguish some keys (e.g. `ctrl+i` == `tab`, `ctrl+m` == `enter`, plain `esc` vs alt combos). Bubble Tea v2 requests basic **key disambiguation** by default from terminals that support the modern keyboard protocols, and tells you what it got via a `tea.KeyboardEnhancementsMsg` (Unit 9), so `shift+enter` or `ctrl+i` vs `tab` *can* work there. But many terminals (and tmux/SSH setups) don't support it, so **always design shortcuts to work on the lowest common denominator**: letters, arrows, enter, esc, ctrl+letter.

### Window size
The runtime sends a `tea.WindowSizeMsg{Width, Height}` once at start and on every resize. Store it:
```go
case tea.WindowSizeMsg:
    m.w, m.h = msg.Width, msg.Height
```
Until it arrives, `m.w`/`m.h` are zero, so guard your `View` (`if m.w == 0 { return tea.NewView("") }`, Unit 6). If you ever need the size again on demand, `tea.RequestWindowSize` is a `func() Msg`, so it is itself a valid `Cmd`: `return m, tea.RequestWindowSize` triggers a fresh `WindowSizeMsg`. Usually unnecessary, since the runtime sends one at startup and on every resize.

### Quit hygiene
Always handle `ctrl+c` explicitly: raw mode means the terminal doesn't turn it into SIGINT for you.

## Worked examples

### Straightforward: navigation keys
```go
case tea.KeyPressMsg:
    switch msg.String() {
    case "up", "k":   m.cursor = max(0, m.cursor-1)
    case "down", "j": m.cursor = min(len(m.items)-1, m.cursor+1)
    }
```
Supporting both arrows and vim keys costs nothing: list them in one `case`.

### Typical: a tiny text field
```go
case tea.KeyPressMsg:
    switch msg.String() {
    case "enter":     return m.submit()
    case "backspace": if n := len(m.text); n > 0 { m.text = m.text[:n-1] }
    default:          m.text += msg.Key().Text
    }
```
Bug hiding here: slicing bytes breaks multibyte characters (é, 日本語). Convert to `[]rune` for real editors, or better, use the Bubbles `textinput` (Unit 7). The `default` branch relies on `Text` being empty for non-typing keys like `up`.

### Edge case: press vs. release and repeated keys
```go
case tea.KeyMsg:               // interface: matches BOTH press and release
    m.count++                  // BUG on enhanced terminals: increments twice per tap
```
On terminals supporting event types (and only if you requested them), a single tap produces a press **and** a release, so a handler on `tea.KeyMsg` fires twice. Fix: match `tea.KeyPressMsg` for actions, or type-switch inside `KeyMsg`. Also, a held key yields repeated presses with `Key().IsRepeat == true`.

## Common misconceptions
- **"`tea.KeyMsg` is what I match, like v1."** In v2 it's an interface, so prefer `KeyPressMsg`.
- **"Space is `" "`."** It's `"space"` from `String()`; `Text` is `" "`.
- **"`msg.Type`/`msg.Runes` exist."** They became `Code` and `Text` (a string).
- **"Uppercase letters need shift detection."** `String()` already returns `"A"` for shift+a.
- **"I can bind ctrl+i and tab differently."** Usually not; assume they're identical.
- **"`WindowSizeMsg` arrives before `Init`."** It arrives *after*, as a normal message; handle the zero-size first frame.

## Check your understanding
1. What does `msg.String()` return for shift+a? For the up arrow? *("A", "up")*
2. Which field would you append to a text buffer, `Code` or `Text`? *(Text; empty for non-typing keys.)*
3. Why might a `case tea.KeyMsg:` handler double-fire on some terminals? *(It matches both press and release.)*
4. Give two messages your own app might define. *(Any: `tickMsg`, `savedMsg`, `errMsg`…)*

## Summary
- Any type can be a message; define your own with distinct types.
- Use `tea.KeyPressMsg` + `msg.String()` for shortcuts; `msg.Key().Text` for typing.
- Space is `"space"`; modifiers prefix as `ctrl+`, `alt+`, `shift+`.
- `WindowSizeMsg` reports terminal size initially and on resize; guard the zero case.
- `KeyMsg` (interface) covers press *and* release.
