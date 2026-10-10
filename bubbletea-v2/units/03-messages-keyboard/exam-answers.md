# Unit 3 Exam — Answers

**E1:**
```go
case tea.KeyPressMsg:
    switch msg.String() {
    case "up", "k":   m.cursor = max(0, m.cursor-1)
    case "down", "j": m.cursor = min(len(m.items)-1, m.cursor+1)
    case "g":         m.cursor = 0
    case "G":         m.cursor = max(0, len(m.items)-1)
    }
```
Note `G` (uppercase) is a separate string from `g`; the `max(0, …)` handles an empty list.

**E2:** (a) `"a"`, Text `"a"`. (b) `"A"`, Text `"A"` (keystroke `"shift+a"`). (c) `"enter"`, Text `""`. (d) `"ctrl+c"`, Text `""`. (e) `"space"`, Text `" "`.

**E3:** `KeyMsg` matches press and release, so on release-enabled terminals `submitted` increments twice per tap. Fix: `case tea.KeyPressMsg:`.

**E4:**
```go
type tickMsg time.Time
type errMsg struct{ err error }
switch msg := msg.(type) {
case tickMsg:            // ...
case errMsg:             m.err = msg.err
case tea.WindowSizeMsg:  m.w, m.h = msg.Width, msg.Height
}
```

**E5:** (a) `String()`: simple, readable, matches `"q"`/`"ctrl+c"`. (b) `Key().Text`: only produces typed text, empty for navigation keys. (c) Field matching (`Code == tea.KeyEnter && Mod == tea.ModShift`) or `String()=="shift+enter"`; note that this requires a terminal with enhanced keyboard protocol, otherwise indistinguishable (so provide an alternative).

**E6:** Slicing removes the last *byte*; `é` is two bytes in UTF-8, leaving an invalid half-character. Fix: `r := []rune(m.text); m.text = string(r[:len(r)-1])` (after a length check). Even that ignores grapheme clusters; prefer `textinput`.

**E7:**
```go
case tea.KeyPressMsg:
    if m.showHelp { m.showHelp = false; return m, nil }
    if msg.String() == "?" { m.showHelp = true; return m, nil }
    // ...normal keys
```
Order matters: the help check comes first so the "swallow" behavior wins.

**E8:** Not safe: on legacy terminals both send byte 0x09, so the app can't tell them apart. Only enhanced-protocol terminals distinguish them; choose different bindings (e.g. `>`/`ctrl+]`) or provide fallbacks.

**E9:**
```go
type connectedMsg struct{}
type disconnectedMsg struct{ reason error }
type incomingMsg struct{ from, text string; at time.Time }
type sendFailedMsg struct{ id string; err error }
```
Good answers use distinct types with payloads, not stringly-typed messages.

**E10:** `m.w` is 0 → `w-2 = -2`, `strings.Repeat` with a negative count **panics**. Fix: `if m.w == 0 { return tea.NewView("") }` at the top of `View`, or `strings.Repeat("─", max(0, w-2))`.
