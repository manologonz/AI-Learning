# Unit 1: Go Survival Kit for TUI Work

## Objectives

By the end of this unit you can:

- Define a struct and attach methods to it, and explain value vs pointer receivers.
- Explain what an interface is and why Go types satisfy them implicitly.
- Use a type switch to branch on the concrete type inside an interface value.
- Write and explain a closure, including one returned from a function.
- Create a module, add a dependency with `go get`, and run a program.
- Read `(value, error)` return pairs and handle `error` idiomatically.

## Prerequisites

None. You only need to have programmed in some language before.

## Why this unit exists

Bubble Tea is small: a handful of interfaces and functions. But all of it leans on five Go ideas: **structs + methods, interfaces, type switches, closures, and multiple returns**. If those are shaky, the framework feels like magic. If they're solid, Bubble Tea feels obvious.

## Core content

### 1. Modules and running code

A Go project is a *module*. You create one with:

```bash
mkdir hello && cd hello
go mod init example.com/hello     # creates go.mod
```

Code lives in `.go` files starting with `package main`; the entry point is `func main()`:

```go
package main

import "fmt"

func main() {
    fmt.Println("hello")
}
```

`go run .` compiles and runs.

> **Version note:** Bubble Tea v2 requires **Go 1.25 or newer**. If your installed Go is older (e.g. 1.24), Go 1.21+ will auto-download the newer toolchain the first time you build a module whose `go.mod` says `go 1.25.0`. Check with `go version` afterwards. `go get <path>` adds a dependency to `go.mod`; `go mod tidy` cleans up. Bubble Tea v2's path is `charm.land/bubbletea/v2`. You'll use it from Unit 2.

### 2. Structs: your data

A struct groups named fields. In Bubble Tea, your entire application state, the **model**, is a struct.

```go
type counter struct {
    count int
    label string
}

c := counter{count: 0, label: "clicks"}
c.count = 5
```

Uppercase names are exported (visible to other packages); lowercase are private to the package. `count` is private; `Count` would be public.

### 3. Methods and receivers

A method is a function attached to a type via a *receiver*:

```go
func (c counter) Describe() string {          // value receiver
    return fmt.Sprintf("%s: %d", c.label, c.count)
}

func (c *counter) Increment() {               // pointer receiver
    c.count++
}
```

- **Value receiver** `(c counter)`: the method gets a **copy**. Changes to `c` don't affect the caller's value.
- **Pointer receiver** `(c *counter)`: the method gets the **address**, so it can mutate the original.

This distinction is crucial for Bubble Tea. v2's `Update` typically *returns a new model* (value style), which is why you'll see the copy semantics used deliberately.

### 4. Interfaces: behavior contracts

An interface is a set of method signatures. **Any type that has those methods satisfies it automatically**. There is no `implements` keyword.

```go
type describer interface {
    Describe() string
}

func show(d describer) {
    fmt.Println(d.Describe())
}

show(counter{count: 3, label: "x"})   // counter has Describe(), so it works
```

Bubble Tea's central type, `tea.Model`, is just an interface. You never "extend" a class; you write the right methods, and the framework accepts your type.

The empty-ish interface `any` (alias of `interface{}`) accepts every type. In Bubble Tea, messages are `tea.Msg`, which is `any`: **anything can be a message**.

### 5. Type switches: "what did I just receive?"

Since a message can be any type, `Update` needs to ask what it actually is:

```go
func classify(v any) string {
    switch x := v.(type) {
    case int:
        return fmt.Sprintf("int %d", x)
    case string:
        return "string " + x
    case error:
        return "error: " + x.Error()
    default:
        return "unknown"
    }
}
```

Inside each `case`, `x` has the **concrete type** for that branch. This one pattern will be the skeleton of every `Update` function you write.

### 6. Multiple returns and errors

Go has no exceptions for ordinary failure. Functions return an `error` as the last value:

```go
n, err := strconv.Atoi("42")
if err != nil {
    return err
}
```

`error` is itself an interface (`Error() string`). `nil` means "no error". In Bubble Tea, errors from async work are turned into messages (Unit 4).

### 7. Functions are values; closures

Functions can be stored in variables, passed, and returned. A **closure** is a function that captures variables from where it was created:

```go
func makeGreeter(name string) func() string {
    return func() string {
        return "hi " + name     // captures `name`
    }
}

g := makeGreeter("Ana")
fmt.Println(g())   // hi Ana
```

Bubble Tea's `tea.Cmd` is literally `func() tea.Msg`. A command is a function the runtime calls later; closures let it carry the data it needs.

### Worked examples

#### **Straightforward: struct + value method**

```go
type point struct{ x, y int }

func (p point) Sum() int { return p.x + p.y }

p := point{2, 3}
fmt.Println(p.Sum())   // 5
```

Step by step: `p` is a `point`; `p.Sum()` passes a copy of `p` as the receiver; the method adds fields; returns 5.

#### **Typical: pointer receiver mutation, and the copy trap**

```go
type point struct{ x, y int }

func (p point) MoveValue()    { p.x++ }   // mutates a copy
func (p *point) MovePointer() { p.x++ }   // mutates the original

p := point{0, 0}
p.MoveValue()
fmt.Println(p.x)   // 0, nothing changed
p.MovePointer()
fmt.Println(p.x)   // 1
```

Why: `MoveValue` received its own copy of `p`, incremented that, and discarded it. `p.MovePointer()` works even though `p` isn't a pointer: Go auto-takes `&p` for addressable variables.

#### **Edge case: closures capturing a loop variable, and type-switch fallthrough of expectations**

```go
var fns []func() int
for i := 0; i < 3; i++ {
    fns = append(fns, func() int { return i })
}
for _, f := range fns {
    fmt.Print(f(), " ")
}
```

Since Go 1.22, each iteration has its **own** `i`, so this prints `0 1 2`. (Before 1.22 it printed `3 3 3`. You'll see the old gotcha in older blog posts.) Your Go 1.24 uses the new semantics.

Second edge: a type switch on a *pointer vs value*. If a message is sent as `myMsg{}` but you write `case *myMsg:`, it **won't match**. The dynamic type is `myMsg`, not `*myMsg`. Bubble Tea messages are sent as values, so match values.

### Common misconceptions

- **"Interfaces must be declared as implemented."** No; satisfaction is implicit and structural.
- **"A method with a value receiver can change my struct."** It changes a copy only.
- **"`any` means untyped."** Values still have a concrete type at runtime; a type switch or assertion recovers it.
- **"`err` being ignored is fine because Go doesn't crash."** An ignored error silently continues with a zero value. Always check.
- **"Capitalization is style."** It is the visibility rule.
- **"A closure copies the captured variable."** It captures the variable itself (shared), not a snapshot.

### Check your understanding

1. Given `func (c counter) Reset() { c.count = 0 }`, does `c.Reset()` zero the caller's counter? *(No. Value receiver mutates a copy.)*
2. If type `T` has methods `A()` and `B()`, does it satisfy `interface{ A() }`? *(Yes; extra methods are fine.)*
3. What does `v.(type)` do, and where is it legal? *(It extracts the dynamic type; only legal inside a `switch`.)*
4. What type is `tea.Cmd`, conceptually? *(A function returning a message: `func() Msg`.)*

## Summary

- A module is a project; `go mod init`, `go get`, `go run .`.
- Structs hold state; methods attach behavior; receivers decide copy vs mutate.
- Interfaces are implicit contracts; `any` accepts everything.
- Type switches recover the concrete type; they're the skeleton of `Update`.
- Errors are returned values; closures capture variables and power `tea.Cmd`.
