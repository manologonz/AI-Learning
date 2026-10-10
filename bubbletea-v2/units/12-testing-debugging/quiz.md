# Unit 12 Quiz

**Q1 (easy).** The v2 teatest import path is:
A) `github.com/charmbracelet/x/exp/teatest`  B) `github.com/charmbracelet/x/exp/teatest/v2`  C) `charm.land/teatest/v2`  D) `charm.land/bubbletea/v2/teatest`

**Q2 (easy).** Where should debug output go?
A) `fmt.Println`  B) A file via `tea.LogToFile`  C) stderr in the TUI  D) `View`

**Q3 (medium).** How do you check that `Update` requested a quit?
A) Compare the model  B) Call the returned `Cmd` and assert it returns `tea.QuitMsg`  C) Read stdout  D) You can't

**Q4 (medium).** What should replace `time.Sleep` when waiting for output in teatest?
A) `runtime.Gosched`  B) `teatest.WaitFor` with a condition and timeout  C) Larger sleeps  D) Nothing

**Q5 (medium).** Why can a teatest test hang at `FinalModel`?
A) The program never quit  B) The terminal is too small  C) Go modules  D) It always hangs

**Q6 (medium).** Short answer: how do you test a command that does an HTTP call?

**Q7 (hard).** Short answer: `cmd := tea.Batch(a, b); msg := cmd()` — what is `msg` and how do you assert on `a`'s result?

**Q8 (medium).** Which env var makes Bubble Tea write a panic log file?
A) `TEA_LOG`  B) `TEA_DEBUG=true`  C) `DEBUG`  D) `BUBBLE_TEA_PANIC`

**Q9 (hard, apply).** Write a table-driven test for `Update` where `"x"` increments and `"z"` resets `count` to 0.

**Q10 (medium).** Name two causes of flaky golden-file tests.
