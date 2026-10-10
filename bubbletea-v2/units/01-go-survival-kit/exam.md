# Unit 1 Exam (summative)

**E1 (medium, Apply).** What does this print? Explain each step.

```go
type stack struct{ items []int }
func (s *stack) Push(n int) { s.items = append(s.items, n) }
func (s stack) Len() int    { return len(s.items) }
func (s stack) Clear()      { s.items = nil }

func main() {
    s := stack{}
    s.Push(1)
    s.Push(2)
    s.Clear()
    fmt.Println(s.Len())
}
```

**E2 (medium, Analyze).** Given the following, say which of the four `var _ Shape = ...` lines compile, and why.

```go
type Shape interface{ Area() float64 }
type Rect struct{ w, h float64 }
type Circle struct{ r float64 }
func (r Rect) Area() float64    { return r.w * r.h }
func (c *Circle) Area() float64 { return 3.14 * c.r * c.r }

var _ Shape = Rect{}
var _ Shape = &Rect{}
var _ Shape = Circle{}
var _ Shape = &Circle{}
```

**E3 (medium, Understand).** In Bubble Tea, `tea.Msg` is `any`. What does that buy the framework, what does it cost you at compile time, and which Go construct restores the ability to handle each message kind safely?

**E4 (hard, Analyze).** This function has two compile errors. Find both and fix them.

```go
func label(v any) string {
    switch x := v.(type) {
    case int, int64:
        return "n=" + strconv.Itoa(x)
    case string:
        return x
    }
}
```

**E5 (medium, Apply).** Write `func counterFrom(start int) func() int` that returns a closure yielding `start+1`, `start+2`, ... on successive calls. Then predict the output of:

```go
a := counterFrom(10)
b := counterFrom(10)
fmt.Println(a(), a(), b())
```

**E6 (hard, Analyze).** Predict the output of each program under Go 1.22+ and explain the difference.

```go
// Program 1
var fs []func() string
for _, name := range []string{"a", "b", "c"} {
    fs = append(fs, func() string { return name })
}
for _, f := range fs { fmt.Print(f()) }

// Program 2
var fs []func() string
var name string
for _, name = range []string{"a", "b", "c"} {
    fs = append(fs, func() string { return name })
}
for _, f := range fs { fmt.Print(f()) }
```

**E7 (medium, Evaluate).** A teammate writes:

```go
f, _ := os.Open(path)
data, _ := io.ReadAll(f)
```

Give one concrete failure this hides and rewrite it to handle errors idiomatically, returning a useful wrapped error.

**E8 (hard, Apply).** Write `func parseAge(s string) (int, error)` that returns an error if `s` is not an integer or is negative. Then show a caller that prints either the age or the error message.

**E9 (hard, Evaluate).** Choose value or pointer receiver for each type and justify: (a) `type Point struct{ X, Y int }` that is never modified after creation; (b) `type Counter struct{ n int }` with an `Increment()` method; (c) `type SafeMap struct{ mu sync.Mutex; m map[string]int }` with `Set` and `Get`. Then state why mixing receiver kinds on one type is discouraged.

**E10 (hard, Create).** Without Bubble Tea, design a tiny message handler: message types `addMsg{n int}` and `resetMsg{}`; a `model` with `total int` and `history []int`; and `func (m model) update(msg any) model` using a type switch, returning a *new* model. `addMsg` adds `n` to `total` and records it in `history`; `resetMsg` clears both. Show two chained calls and the final state. Bonus: explain how a careless `append` to `history` could corrupt an earlier model value, and how to avoid it.
