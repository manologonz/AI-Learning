# Unit 1 Quiz (formative)

**Q1 (easy).** What does `go mod init example.com/app` create?
A) A `main.go` file  B) A `go.mod` file  C) A `vendor/` directory  D) A compiled binary

**Q2 (easy).** Which identifier is exported from its package?
A) `count`  B) `_count`  C) `Count`  D) `count_`

**Q3 (medium).** Given `func (p point) Move() { p.x++ }` and `p := point{}`, after `p.Move()` what is `p.x`?
A) 1  B) 0  C) A compile error  D) Undefined behavior

**Q4 (medium).** How does a type satisfy an interface in Go?
A) Using `implements`  B) By embedding the interface  C) By having all its methods  D) By registering it with `init()`

**Q5 (medium).** Short answer: in one sentence, why does `switch m := msg.(type)` need `msg` to be an interface type?

**Q6 (medium).** What does this print?

```go
func f() func() int { n := 0; return func() int { n++; return n } }
g := f(); g(); g()
fmt.Println(g())
```

A) 1  B) 3  C) 0  D) 2

**Q7 (medium).** `x, err := strconv.Atoi("abc")`. What are `x` and `err`?
A) x=0, err=non-nil  B) x=nil, err=non-nil  C) Program panics  D) x=0, err=nil

**Q8 (hard).** Application: write a `type temp float64` with a method `Fahrenheit() float64` (value receiver) converting Celsius.

**Q9 (hard).** Given

```go
type stopMsg struct{}
func handle(m any) string {
    switch m.(type) {
    case *stopMsg: return "stop"
    default: return "other"
    }
}
```

what does `handle(stopMsg{})` return, and why?

**Q10 (medium).** Short answer: what is `tea.Cmd` conceptually, and why is a closure a natural way to build one?
