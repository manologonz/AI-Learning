# Unit 8 Quiz

**Q1 (easy).** How should a child inform its parent of an event?
A) Call a parent method via pointer  B) Return a command that produces a message  C) Write to a global variable  D) `os.Exit`

**Q2 (easy).** In a router's `Update`, where should global `ctrl+c` handling go?
A) In every child  B) Before delegating to children  C) Only in `View`  D) In `Init`

**Q3 (medium).** Which messages should typically go to *all* children?
A) Key presses  B) Ticks/blink/size/async results  C) Only mouse events  D) None

**Q4 (medium).** What's a downside of `tea.Model`-typed children vs concrete types?
A) They can't render  B) Type assertions/indirection  C) They can't receive messages  D) They can't return commands

**Q5 (medium).** Short answer: who should own the terminal size, parent or child, and why?

**Q6 (medium).** A form shows the previous input each time you revisit it. Likely cause?

**Q7 (hard).** Short answer: an HTTP reply for screen A arrives after the user navigated to screen B. Describe two acceptable behaviors and how to implement one.

**Q8 (hard, apply).** Write a `navMsg` type and the child command that emits it to go to `settingsScreen`.

**Q9 (medium).** What's wrong with `m.form.Update(msg); return m, nil`?

**Q10 (medium).** Why is having `cursor` in both parent and child risky?
