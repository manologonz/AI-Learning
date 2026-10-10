# Unit 10 Quiz

**Q1 (easy).** Which goroutine runs `Update`?
A) A new goroutine per message  B) One event-loop goroutine  C) The command's goroutine  D) `main` before `Run`

**Q2 (easy).** How does an external goroutine deliver an event to the UI?
A) Mutate the model  B) `p.Send(msg)`  C) Call `Update` directly  D) Write to stdout

**Q3 (medium).** What does `Send` do after the program exited?
A) Panics  B) Blocks forever  C) No-op  D) Returns an error

**Q4 (medium).** A command waiting on a channel returns one message. To keep listening you:
A) Loop inside the command  B) Return the same wait command again from `Update`  C) Use `tea.Batch`  D) Nothing; it loops automatically

**Q5 (medium).** Which stops a program immediately without the final render?
A) `p.Quit()`  B) `tea.Quit`  C) `p.Kill()`  D) `p.Wait()`

**Q6 (medium).** Short answer: what does `tea.WithFilter` let you do?

**Q7 (hard).** Short answer: your producer sends 10,000 events/second and the UI lags. Why, and what do you change?

**Q8 (medium).** True/False: quitting the program cancels running commands.

**Q9 (hard, apply).** Write `waitForLine(ch <-chan string) tea.Cmd` that returns `lineMsg` or `streamEndedMsg` if the channel closes.

**Q10 (medium).** Why can reading `m.count` from a debug goroutine be a bug?
