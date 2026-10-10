# Unit 8: Composing Models — Parent/Child & Multi-Screen Apps

> Runnable code: `examples/u08-screens/main.go` (menu → form → done, with a router model).

## Objectives
By the end of this unit you can:
- Structure an app as a tree of models: a parent that owns children and routes messages.
- Design child→parent communication with messages (not shared pointers).
- Implement multi-screen navigation as a state machine (enum + router).
- Manage focus so only one child receives key input, while everyone gets non-key messages.
- Distribute window size to children and decide who owns layout.

## Prerequisites
Units 2–7.

## Core content

### Why compose?
A single `Update` with a 500-line switch is unmaintainable. Bubbles (Unit 7) already showed the alternative: a model that embeds child models and forwards messages. Apply that pattern to *your own* screens and widgets: each is a small struct with its own state, `Update`, `View`.

### Two shapes for a child
1. **Concrete-typed child** (simplest; used in `examples/u08-screens`):
   ```go
   type menu struct{ cursor int }
   func (m menu) Update(msg tea.Msg) (menu, tea.Cmd) { … }   // returns menu, not tea.Model
   func (m menu) View() string { … }                          // string: parent composes them
   ```
   Parent stores `menu` by value; forwarding is `m.menu, cmd = m.menu.Update(msg)` with no type assertions.
2. **`tea.Model`-typed child** (when the child set is dynamic, e.g. a stack of screens): store `tea.Model`, and `Update` returns `tea.Model`:
   ```go
   type router struct{ current tea.Model }
   m.current, cmd = m.current.Update(msg)
   ```
   Cost: interface indirection and occasional type assertions. Benefit: swapping screens is `m.current = newScreen()`.

Choose (1) unless you need polymorphism.

### The router pattern
```go
type screen int
const ( menuScreen screen = iota; formScreen; doneScreen )

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // 1. global concerns (quit, window size) handled first
    if k, ok := msg.(tea.KeyPressMsg); ok && k.String() == "ctrl+c" { return m, tea.Quit }
    // 2. messages addressed to the router (navigation, results from children)
    switch msg := msg.(type) {
    case navMsg:    m.current = msg.to; return m, nil
    case submitMsg: m.name, m.current = msg.name, doneScreen; return m, nil
    }
    // 3. everything else to the ACTIVE child only
    var cmd tea.Cmd
    switch m.current {
    case menuScreen: m.menu, cmd = m.menu.Update(msg)
    case formScreen: m.form, cmd = m.form.Update(msg)
    }
    return m, cmd
}
```
Order matters: global → router-level → delegated.

### Child → parent: messages, not shared state
A child can't call the parent. It *returns a command that yields a message*, and the parent handles it:
```go
// child:  on enter
return f, func() tea.Msg { return submitMsg{name: f.text} }
```
Now the child has no knowledge of the parent (reusable, testable), and the parent decides what "submit" means. Avoid passing pointers to parent state into children; that reintroduces the coupling MVU removes. Note that the message goes through the runtime loop, arriving at the parent's `Update` next; think of it as an event.

Parent → child is direct (call `Update`, set a field, or pass constructor args).

### Focus
Only one interactive child should get keyboard input. Model it explicitly:
```go
type focus int
const ( focusList focus = iota; focusFilter )
// key messages → only the focused child
// non-key messages (ticks, blink, size, async results) → everyone who cares
```
Anti-pattern: forwarding *every* message to *every* child. Two text inputs would both type. Conversely, forwarding **only** key messages to the focused child and **nothing else** starves animations (spinner ticks, cursor blink). Split by message kind.

### Window size and layout ownership
Rule of thumb: **the parent owns layout**. It receives `WindowSizeMsg`, computes each child's rectangle, and tells children their size (`child.SetSize(w,h)` or a `resizeMsg`). Children render into the size they're given and don't read the terminal. This makes children reusable in different layouts and easy to test at fixed sizes.

### Lifecycle: initializing a screen on entry
When switching screens you often need to (re)initialize state and start commands:
```go
case navMsg:
    m.current = msg.to
    if msg.to == formScreen {
        m.form = newForm()                    // fresh state
        return m, m.form.Init()               // and start its commands (e.g. cursor blink)
    }
```
Forgetting this is why "the form still shows what I typed last time" and "the cursor doesn't blink on the second visit".

### Keeping messages unambiguous
Different screens can reuse key names for different meanings, which is fine *because only the active screen receives them*. But async results (an HTTP reply for screen A) may arrive while B is active. Route those by **type**, not by which screen is showing, or carry an ID (Unit 4) and ignore mismatches. A reply must never be dropped just because the user navigated away *unless* you intend that.

## Worked examples

### Straightforward: menu → done
Trace `examples/u08-screens`: `enter` on "Enter name" → menu returns a `navMsg{formScreen}` command → router receives it, sets `current` → next `View` renders the form.

### Typical: list + detail panes with focus
```go
type model struct {
    list   list.Model
    detail viewport.Model
    focus  focus
    w, h   int
}
case tea.WindowSizeMsg:
    m.w, m.h = msg.Width, msg.Height
    m.list.SetSize(m.w/3, m.h-1)
    m.detail.SetWidth(m.w - m.w/3)
    m.detail.SetHeight(m.h-1)
case tea.KeyPressMsg:
    if msg.String() == "tab" { m.focus = 1 - m.focus; return m, nil }
    // route key to focused child only
    if m.focus == focusList { m.list, cmd = m.list.Update(msg) } else { m.detail, cmd = m.detail.Update(msg) }
    return m, cmd
default: // non-key: forward to both
```
After the list selection changes, the parent reads `m.list.SelectedItem()` and calls `m.detail.SetContent(...)`. *Parent-orchestrated* data flow.

### Edge case: value-copy pitfalls in nested models
```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    m.form.Update(msg)          // BUG: returned model discarded
    return m, nil
}
```
Also: a child that returns `tea.Model` from a concrete-typed slot needs a type assertion (`m.child = next.(child)`), and a failing assertion **panics**. When the child type is concrete, have its `Update` return the concrete type to avoid this. And beware **two sources of truth**: storing `cursor` in both parent and child, then updating only one. Each fact lives in exactly one model; others derive from it.

## Common misconceptions
- **"Children can hold a pointer to the parent to call it."** It couples them and hides data flow; send messages.
- **"Forward every message to every child."** Causes duplicate input handling; split by key vs non-key and focus.
- **"Only key messages matter to children."** Ticks, blinks, size, and async results are messages too.
- **"A screen's state persists on its own."** Reset/initialize on entry, or you'll see stale state.
- **"The child should read the terminal size."** The parent should hand it its size.
- **"Interface-typed children are always better."** Concrete types are simpler and safer.

## Check your understanding
1. How does a child tell its parent "user submitted"? *(Return a command producing a message.)*
2. Why forward non-key messages even to unfocused children? *(Animations/async results.)*
3. In the router, why handle nav messages before delegating? *(They're addressed to the router.)*
4. Who computes child sizes? *(The parent.)*

## Summary
- Compose with small child models; parent owns routing, focus, and layout.
- Child→parent via messages; parent→child via calls/fields.
- Route keys to the focused child; forward non-key messages broadly.
- Initialize a screen when entering it; keep one source of truth per fact.
