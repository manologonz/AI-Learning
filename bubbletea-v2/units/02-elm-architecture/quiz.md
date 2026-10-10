# Unit 2 Quiz

**Q1 (easy).** Which is the correct v2 `Model` interface?
A) `Init() (Model, Cmd); Update(Msg) (Model, Cmd); View() string`
B) `Init() Cmd; Update(Msg) (Model, Cmd); View() View`
C) `Init() Cmd; Update(Msg) Model; View() View`
D) `Init(); Update(Msg); View() string`

**Q2 (easy).** Who calls `View()`?
A) Your `Update`  B) `main`  C) The Bubble Tea runtime, after each `Update`  D) A command

**Q3 (medium).** `Update` runs `m.count++` then `return model{}, nil`. Result?
A) count increments  B) count stays 0/resets  C) compile error  D) panic

**Q4 (medium).** What is `tea.Quit`?
A) A function that calls `os.Exit(0)`  B) A `Cmd` that leads to a `QuitMsg`  C) A key constant  D) A model method

**Q5 (medium).** In v2, which case matches the space bar via `msg.String()`?
A) `" "`  B) `"space"`  C) `"spacebar"`  D) `"ctrl+space"`

**Q6 (medium).** Short answer: why does the model use a value receiver yet changes persist?

**Q7 (medium).** Short answer: name what `p.Run()` returns.

**Q8 (hard).** A teammate stores a `map[string]bool` in the model and mutates it inside `Update`. State an advantage they get and a risk they take.

**Q9 (medium).** Why is `fmt.Println("debug")` inside `Update` a bad idea?

**Q10 (hard, apply).** Write the `Update` case (inside a `tea.KeyPressMsg` switch) that toggles a `bool` field `paused` on `"p"` and quits on `"q"`.
