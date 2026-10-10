# Unit 1 Exam: Answers

**E1: 2.** `Push` has a pointer receiver, so `s.Push(1)` (Go auto-takes `&s`) appends to the real `s.items`; after two pushes it holds `[1 2]`. `Clear` has a value receiver: it sets `items = nil` on a copy, leaving `s` untouched. `Len()` therefore returns 2.

**E2:**
- `Rect{}` ✓: `Area` has a value receiver, so it is in `Rect`'s method set.
- `&Rect{}` ✓: a pointer's method set includes the value-receiver methods too.
- `Circle{}` ✗: `Area` has a pointer receiver, so only `*Circle` has it in its method set.
- `&Circle{}` ✓.

Rule: `T` has value-receiver methods; `*T` has both value- and pointer-receiver methods.

**E3:** Because `Msg` is `any`, any type (key press, window size, your own struct) can be a message, so the framework needs no registry or base type and you can define new message kinds freely. The cost: the compiler no longer knows what a message is, so a wrong assumption is a runtime issue, not a compile error. A **type switch** (`switch m := msg.(type)`) recovers the concrete type in each branch, with a `default` for messages you don't care about.

**E4:**
1. In a multi-type case (`case int, int64:`), `x` keeps the interface type `any`, so `strconv.Itoa(x)` doesn't compile (it needs an `int`). Fix: split into `case int:` and `case int64:` (using `strconv.FormatInt(x, 10)`), or format with `fmt.Sprint(x)`.
2. The function has no return on the path where no case matches: "missing return". Fix: add a `default:` returning a value, or a final `return`.

```go
func label(v any) string {
    switch x := v.(type) {
    case int:
        return "n=" + strconv.Itoa(x)
    case int64:
        return "n=" + strconv.FormatInt(x, 10)
    case string:
        return x
    default:
        return "?"
    }
}
```

**E5:**
```go
func counterFrom(start int) func() int {
    return func() int {
        start++
        return start
    }
}
```
Output: `11 12 11`. Each call to `counterFrom` creates its own `start` variable; the closure captures that variable, so `a` and `b` have independent state, and `a` mutates its own between calls.

**E6:** Program 1 prints `abc`; Program 2 prints `ccc`. In Program 1, `:=` in the range clause gives each iteration its own `name` (Go 1.22+), so each closure captures a different variable. In Program 2, `name` is declared outside and the loop merely assigns to it; all three closures share that one variable, which holds `"c"` after the loop ends. Closures capture variables, not values.

**E7:** If `path` doesn't exist, `os.Open` fails and `f` is nil; the code then continues and `data` is empty, indistinguishable from an empty file (a missing file is silently treated as "no content"; a read failure mid-way is likewise lost). Fix:

```go
f, err := os.Open(path)
if err != nil {
    return nil, fmt.Errorf("open %s: %w", path, err)
}
defer f.Close()
data, err := io.ReadAll(f)
if err != nil {
    return nil, fmt.Errorf("read %s: %w", path, err)
}
```

**E8:**
```go
func parseAge(s string) (int, error) {
    n, err := strconv.Atoi(s)
    if err != nil {
        return 0, fmt.Errorf("invalid age %q: %w", s, err)
    }
    if n < 0 {
        return 0, fmt.Errorf("age must be non-negative, got %d", n)
    }
    return n, nil
}

age, err := parseAge(input)
if err != nil {
    fmt.Println("error:", err)
    return
}
fmt.Println("age:", age)
```

**E9:** (a) Value receiver: small, immutable, cheap to copy, safe to share. (b) Pointer receiver: `Increment` must mutate the caller's value; a value receiver would increment a copy. (c) Pointer receiver is required: copying a struct containing a `sync.Mutex` copies the lock (a bug, flagged by `go vet`), and the methods mutate shared state. Mixing receivers on one type is discouraged because the method sets differ (`T` vs `*T`), which leads to confusing interface-satisfaction failures and to readers being unsure whether methods mutate.

**E10:**
```go
type addMsg struct{ n int }
type resetMsg struct{}

type model struct {
    total   int
    history []int
}

func (m model) update(msg any) model {
    switch msg := msg.(type) {
    case addMsg:
        m.total += msg.n
        m.history = append(slices.Clone(m.history), msg.n)
    case resetMsg:
        m.total = 0
        m.history = nil
    }
    return m
}

m := model{}.update(addMsg{3}).update(addMsg{4})
// m == model{total: 7, history: []int{3, 4}}
```
The value receiver means `m` is a copy that we modify and return, the same shape as Bubble Tea's `Update`. Bonus: a plain `append(m.history, n)` may write into the shared backing array when it has spare capacity, so two models derived from the same earlier model can overwrite each other's history (and the old model's slice contents can change under it). Cloning first (or `append(m.history[:len(m.history):len(m.history)], n)`) forces a fresh array.
