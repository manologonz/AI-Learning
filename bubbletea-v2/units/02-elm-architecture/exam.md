# Unit 2 Exam

**E1 (Apply).** Write a complete program (imports, model, three methods, main) that shows "Hello, <name>" where name starts as `"world"`, and pressing `n` switches it to `"Go"`, `q` quits. Must compile against v2.

**E2 (Analyze).** Trace: model `{count:0}`, messages in order: `up`, `up`, `down`, `x`, `q` with the counter program. Give the model after each message and what `Run()` returns.

**E3 (Analyze).** Find every v2 error:
```go
func (m model) Init() (tea.Model, tea.Cmd) { return m, nil }
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == " " { m.n++ }
    }
}
func (m model) View() string { return "n" }
```

**E4 (Understand).** Explain the "spreadsheet" analogy for MVU in 3 sentences, naming what corresponds to cells and formulas.

**E5 (Evaluate).** Two designs for a todo app: (A) global variables plus `fmt.Print` redraws; (B) MVU. Give two concrete advantages of B.

**E6 (Analyze).** `m.items = append(m.items, x)` inside `Update` (value receiver). Is the change kept? Is there any aliasing risk? Explain, referencing slice capacity.

**E7 (Apply).** Add bounds so `count` stays within 0..10 in the counter's `Update`.

**E8 (Analyze).** Why is it fine that `Update` mutates `m` (a copy) but wrong for `View` to do so? Give a bug scenario.

**E9 (Evaluate).** A colleague calls `os.Exit(0)` on `q`. List two things that go wrong and the correct alternative.

**E10 (Create).** Design the model struct (fields only) for a tiny stopwatch UI with start/stop/reset and a displayed elapsed time. Which fields are state and which are derived in `View`? (No timing code yet; Unit 4 covers ticks.)
