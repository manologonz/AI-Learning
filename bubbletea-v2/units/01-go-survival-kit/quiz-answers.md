# Unit 1 Quiz — Answers

**Q1: B.** `go mod init` creates `go.mod`. A is wrong: it never writes source. C is wrong: vendoring is `go mod vendor`. D is wrong: building is `go build`.

**Q2: C.** Exported = starts with an uppercase letter. A, B, D start with lowercase or underscore, so they're package-private.

**Q3: B (0).** Value receiver: the method mutates a copy. A is the classic mistaken belief that methods mutate. C: it compiles fine. D: the behavior is fully defined.

**Q4: C.** Structural, implicit satisfaction. A: Go has no `implements`. B: embedding is a different mechanism (composition). D: no registration exists.

**Q5:** A type switch inspects the *dynamic* type stored inside an interface value; a concrete-typed variable has only one possible type, so there is nothing to discriminate (and the compiler rejects it).

**Q6: B (3).** The closure captures `n` by reference; three calls increment it to 3. A treats each call as fresh state. C/D miscount calls.

**Q7: A.** `Atoi` returns the zero value 0 plus a non-nil error. B: `int` can't be `nil`. C: Go returns an error rather than panicking. D: the input isn't a number, so `err` isn't nil.

**Q8:**
```go
type temp float64
func (t temp) Fahrenheit() float64 { return float64(t)*9/5 + 32 }
```
Value receiver is right because nothing is mutated. Key point: conversion `float64(t)` is needed because `temp` and `float64` are distinct types.

**Q9:** `"other"`. The dynamic type is `stopMsg` (value), but the case matches `*stopMsg` (pointer). They're different types, so it falls to `default`. Fix: `case stopMsg:`.

**Q10:** `tea.Cmd` is a function returning a message (`func() tea.Msg`) that the runtime executes later, off the main loop. Closures capture the data the command needs (a URL, an ID, a duration) without a separate struct.
