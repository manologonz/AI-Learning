# Unit 2 Exam — Answers

**E1:**
```go
package main
import ("fmt"; "os"; tea "charm.land/bubbletea/v2")
type model struct{ name string }
func (m model) Init() tea.Cmd { return nil }
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if k, ok := msg.(tea.KeyPressMsg); ok {
        switch k.String() {
        case "n": m.name = "Go"
        case "q", "ctrl+c": return m, tea.Quit
        }
    }
    return m, nil
}
func (m model) View() tea.View { return tea.NewView("Hello, " + m.name + "\n") }
func main() {
    if _, err := tea.NewProgram(model{name: "world"}).Run(); err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
}
```
Grading points: `tea.KeyPressMsg`, `Init() tea.Cmd`, `View() tea.View`, `tea.NewView`, initial `name`, returning `m`.

**E2:** up→1, up→2, down→1, `x`→1 (no case matches; returns `m,nil`), `q`→ returns `tea.Quit`. `Run()` returns `model{count:1}, nil`.

**E3:** (1) `Init` must return only `tea.Cmd`. (2) `tea.KeyMsg` is an interface in v2; use `tea.KeyPressMsg` for presses. (3) space is `"space"`, not `" "`. (4) `Update` has no return statement (missing `return m, nil`). (5) `View` must return `tea.View`, e.g. `tea.NewView("n")`. (6) `m.n` requires the field declared and is fine, but `m.n++` result is never returned.

**E4:** Model = the cells (the only stored data); View = formulas that compute what's displayed from the cells; Update = the only place you edit cells, in response to events. Result: display always consistent with data because it's derived.

**E5:** (1) Screen can't drift out of sync since it's derived from one model. (2) Easier testing: `Update` is a pure-ish function of (model,msg). Also: no scattered redraw calls; single-threaded event handling avoids locks.

**E6:** The new slice header is stored in the copy you return, so the change is kept. Aliasing risk: if capacity allows, `append` writes into the shared backing array, so an old model copy (same slice header, shorter len) still wouldn't *see* the element but the underlying array is shared; with `m.items[i] = v` (index assignment) both copies see the change. Takeaway: don't rely on old-model immutability unless you copy slices.

**E7:**
```go
case "up":   if m.count < 10 { m.count++ }
case "down": if m.count > 0  { m.count-- }
```

**E8:** `Update` is the sanctioned place where state transitions happen, and the result is explicitly returned. `View` is called after each update, possibly more often than you expect (e.g. re-rendering), so mutating there makes output depend on how many times it's called, e.g. a `View` that increments a frame counter would race ahead and produce non-reproducible screens.

**E9:** The terminal stays in raw mode / alt screen / hidden cursor (unusable shell); deferred cleanup and `Run()`'s return value are skipped. Correct: return `tea.Quit` from `Update`.

**E10:**
```go
type model struct {
    running  bool
    elapsed  time.Duration // state: accumulated time
    lastTick time.Time     // state: for accumulating (optional)
}
```
Derived in `View`: the formatted string (`mm:ss.t`), the label "Running"/"Stopped". Don't store formatted strings.
