# Unit 2: The Elm Architecture & Your First Program

> Runnable code for this unit: `examples/u02-counter/main.go` (`go run ./examples/u02-counter`).

## Objectives
By the end of this unit you can:
- Explain the Model → Update → View loop and what triggers each step.
- Write a complete Bubble Tea v2 program from an empty file.
- Explain why `Update` returns a (new) model and a command, and what `tea.Quit` is.
- Explain what breaks if `View` mutates state or `Update` forgets to return the model.
- Describe what `Program.Run()` returns and when it returns.

## Prerequisites
Unit 1 (structs, methods, receivers, interfaces, type switches).

## Core content

### The problem Bubble Tea solves
A terminal app must react to unpredictable events (keys, resizes, timers, network replies) and keep the screen consistent with its state. The classic disaster is scattering "draw this now" calls everywhere so screen and state drift apart.

Bubble Tea uses **The Elm Architecture** (also called MVU: Model-View-Update). Three rules:

1. All state lives in one **model** value.
2. Anything that happens becomes a **message** (`tea.Msg`). Only `Update` reacts to messages.
3. The screen is a pure **function of the model**: `View(model) → what to show`. You never "draw"; you describe.

Analogy: a spreadsheet. Cells (model) change only when you edit; every displayed value is a formula of cells (view). You never manually repaint a formula result.

### The three methods
Any type with these three methods is a `tea.Model` (implicit interface satisfaction from Unit 1):

```go
type Model interface {
    Init() Cmd
    Update(Msg) (Model, Cmd)
    View() View
}
```

- `Init() tea.Cmd`: called once at start. Returns an optional first command (`nil` = none).
- `Update(msg) (tea.Model, tea.Cmd)`: called for each message. Returns the **next** model and an optional command (Unit 4).
- `View() tea.View`: called after every `Update` to produce the screen. In v2 it returns a `tea.View` struct, not a string (v1 returned `string`; Unit 5 explores the struct). For plain text: `tea.NewView("...")`.

> **v2 note:** In Bubble Tea v2, `Init()` returns only a `Cmd`. Some older tutorials show `Init() (Model, Cmd)`. That's wrong for v2.

### The loop
```
      ┌──────────────── message (key, tick, resize, your own…) ◄─┐
      ▼                                                          │
  Update(msg) ──► new model, cmd ──► View(model) ──► screen      │
                          │                                      │
                          └── runtime runs cmd in a goroutine ───┘ (its result becomes a message)
```
The runtime owns the loop. **Your code never calls `Update` or `View` itself.** It's the framework's job. The event loop runs `Update` on a single goroutine, one message at a time, which is why your model needs no locks.

### Your first program, annotated
```go
package main

import (
    "fmt"
    "os"

    tea "charm.land/bubbletea/v2"
)

type model struct{ count int }

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {          // Unit 1: type switch
    case tea.KeyPressMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit          // a Cmd that tells the runtime to exit
        case "up", "k":
            m.count++                   // m is a COPY; we return the changed copy
        case "down", "j":
            m.count--
        }
    }
    return m, nil                       // ALWAYS return the model you want kept
}

func (m model) View() tea.View {
    return tea.NewView(fmt.Sprintf("Count: %d\n\n↑/↓ to change, q to quit\n", m.count))
}

func main() {
    if _, err := tea.NewProgram(model{}).Run(); err != nil {
        fmt.Fprintln(os.Stderr, "error:", err)
        os.Exit(1)
    }
}
```

Setup:
```bash
mkdir counter && cd counter
go mod init example.com/counter
go get charm.land/bubbletea/v2
# paste the code into main.go
go run .
```

### Why value receivers and returning the model
`Update` with a value receiver mutates a **copy** of the model and returns it, exactly the Unit 1 copy semantics. The runtime replaces its stored model with whatever you return. Consequence: **if you change `m` and then `return nil-model`/a different variable, your change is lost.** (You *can* use pointer receivers, but the idiomatic v2 style, and what all the official examples use, is value receivers.)

### What `Run()` returns
`p.Run()` blocks until the program exits, then returns the **final model** and an error. You can inspect the final model after the UI is gone (handy for "what did the user choose?" prompts and for tests, Unit 12).

### What `tea.Quit` really is
`tea.Quit` is a `Cmd` (a function returning a `QuitMsg`). Returning it from `Update` says "run this command"; the runtime executes it, gets `QuitMsg`, and shuts down cleanly, restoring the terminal. Never call `os.Exit` inside a TUI. You'd leave the terminal in raw mode.

### Ctrl+C
In raw mode, Ctrl+C is delivered as a key press (`"ctrl+c"`); *you* decide to quit on it. Forgetting to handle it means the user can't exit (a very common first-program bug).

## Worked examples

### Straightforward: the counter
Trace pressing `↑` twice then `q`:
1. Start: `model{count:0}`; `View` shows "Count: 0".
2. `↑` → `KeyPressMsg` with `String()=="up"` → `Update` copies model, `count=1`, returns it; `View` shows "Count: 1".
3. `↑` → count 2. 
4. `q` → returns `tea.Quit`; loop ends; `Run()` returns `model{count:2}, nil`.

### Typical: a cursor over a list
```go
type model struct {
    choices  []string
    cursor   int
    selected map[int]struct{}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if k, ok := msg.(tea.KeyPressMsg); ok {
        switch k.String() {
        case "up", "k":
            if m.cursor > 0 { m.cursor-- }
        case "down", "j":
            if m.cursor < len(m.choices)-1 { m.cursor++ }
        case "space", "enter":
            if _, ok := m.selected[m.cursor]; ok {
                delete(m.selected, m.cursor)
            } else {
                m.selected[m.cursor] = struct{}{}
            }
        case "q":
            return m, tea.Quit
        }
    }
    return m, nil
}
```
Real programs look like this: state (cursor, selection) in the model; keys change it; `View` loops over `choices` and draws `>` at `m.cursor`. Note the bounds checks: the model must never enter an invalid state.

### Edge case: copies that still share memory
```go
case "space":
    m.selected[m.cursor] = struct{}{}   // mutates the map the OLD model also points to
```
The model is copied, but **maps and slices inside it are shared references**. Mutating them mutates "both" models. In a small app this "works" (and the official examples do it), but it breaks assumptions when you compare old/new models, keep history/undo, or share models across goroutines. Safe alternatives: build a new map/slice on change, or treat the model as the single owner and never retain old copies. The lesson: *copy semantics protect scalar fields, not reference fields.*

Second edge: `case " ":` never matches in v2, because space's `String()` is `"space"` (v1 used `" "`).

## Common misconceptions
- **"I call `View()` when I want to redraw."** You never call it; the runtime does after every `Update`.
- **"`Update` should print to the screen."** It only computes the next model; printing there corrupts the display.
- **"`Init` returns the model."** Not in v2; it returns only a `Cmd`.
- **"Value receiver means the model can't change."** It changes by *being returned*.
- **"`return m, tea.Quit` quits immediately, so later code doesn't matter."** `Quit` is a command; `Update`'s remaining statements run normally; the quit happens after you return.
- **"I can use `fmt.Println` for debugging."** It fights the renderer (Unit 13: log to a file).

## Check your understanding
1. Which of Init/Update/View may run more than once? *(Update and View; Init once.)*
2. You write `m.count++` in `Update` but end with `return model{}, nil`. What does the user see? *(Count never changes; you returned a fresh zero model.)*
3. What does `Run()` return? *(Final model and an error.)*
4. Why can't `Update` safely be called from two goroutines by your own code? *(The runtime guarantees single-threaded Update; calling it yourself would race with the loop; send a message instead.)*

## Summary
- MVU: state in one model; events as messages; screen as a function of the model.
- v2 `Model`: `Init() Cmd`, `Update(Msg) (Model, Cmd)`, `View() View`.
- Value receivers + returning the modified copy is the idiomatic pattern; reference fields (maps, slices) are still shared.
- `tea.Quit` is a command; `Run()` returns the final model.
- Always handle `ctrl+c`; match space as `"space"`.
