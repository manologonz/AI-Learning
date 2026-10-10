# Unit 4 Quiz

**Q1 (easy).** `tea.Cmd` is:
A) `func(Model) Model`  B) `func() Msg`  C) `func(Msg) Cmd`  D) `chan Msg`

**Q2 (easy).** What does returning `nil` as the Cmd mean?
A) Quit  B) Error  C) No side effect  D) Re-run Update

**Q3 (medium).** Which is right for "save file, then quit"?
A) `tea.Batch(save, tea.Quit)`  B) `tea.Sequence(save, tea.Quit)`  C) Call `os.Exit` after saving  D) `tea.Tick(0, quit)`

**Q4 (medium).** After handling a `tickMsg` you `return m, nil`. What happens to the timer?
A) Continues  B) Stops  C) Speeds up  D) Panics

**Q5 (medium).** Where does a command's function body run?
A) On the event loop goroutine  B) In `View`  C) On a separate goroutine managed by the runtime  D) At compile time

**Q6 (medium).** Short answer: why should a command return an `errMsg` instead of just failing silently?

**Q7 (hard).** Short answer: you fire a search on each keystroke; replies arrive out of order. Describe the fix.

**Q8 (medium).** True/False: `tea.Every` and `tea.Tick` are identical.

**Q9 (hard, apply).** Write `func fetchName(id int) tea.Cmd` that returns `nameMsg{id, name}` from a hypothetical `lookup(id) (string, error)` or an `errMsg` on failure.

**Q10 (medium).** What's wrong with `time.Sleep(2*time.Second)` in `Update`?
