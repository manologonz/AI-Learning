# Unit 15 Quiz (design check before you build)

**Q1.** Which package must **not** import Bubble Tea in the suggested architecture?
A) `ui`  B) `core`  C) `cmd/taskmate`  D) none

**Q2.** Why do toast-expiry messages carry an ID?
A) To sort them  B) To ignore stale timers when a newer status replaced the old one  C) For logging  D) Because Tick requires it

**Q3.** The list's filter mode is active and the user types `q`. What must the app do?
A) Quit  B) Let the list receive it (skip the global quit binding while filtering)  C) Ignore  D) Clear the filter

**Q4.** Best way to persist tasks safely?
A) `os.WriteFile` directly  B) Write to a temp file then `os.Rename`  C) Write in `View`  D) Keep only in memory

**Q5.** Give the data states for the task list screen.

**Q6.** Which sequence guarantees save-before-quit?
A) `tea.Batch(save, tea.Quit)`  B) `tea.Sequence(save, tea.Quit)`  C) `os.Exit`  D) `tea.Tick`

**Q7.** Why keep the selected task's ID instead of its index?

**Q8.** Where does the width math for the detail pane come from in v2 (border/padding)?
