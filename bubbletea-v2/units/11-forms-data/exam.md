# Unit 11 Exam

**E1 (Apply).** Add a third field "age" (integer 18–120) to `examples/u11-form`: describe validator, focus handling changes, submit checks, and table column.

**E2 (Analyze).** After a refresh the highlighted row jumps to a different record. Diagnose and fix in code.

**E3 (Evaluate).** For a tool that asks 4 questions at install time, and for an in-app "edit profile" screen with a live avatar preview: choose `huh` standalone, `huh` embedded, or hand-rolled, and justify.

**E4 (Apply).** Implement the explicit state enum for a data screen, including transitions in `Update` and the `View` switch.

**E5 (Analyze).** Typing in a 50k-row filter box lags. List three causes and improvements.

**E6 (Apply).** Write a function `visible(all []Record, query string, sortBy string) []Record` (case-insensitive contains on `Name`; sort stable by field) and describe when it's called.

**E7 (Analyze).** What can go wrong when `huh.Form.Update` returns a `tea.Model` and you write `m.form = form.(*huh.Form)` unguarded? Should you guard, and how?

**E8 (Evaluate).** Live validation vs on-blur validation for: email, password strength, username availability, max length.

**E9 (Analyze).** "loading" is true and `err` is set after a retry. Show how a state enum avoids this.

**E10 (Create).** Design a "bulk edit" UI: multi-select rows in a table, apply an action, show progress and per-row failures. Specify state, messages, and how partial failure is represented.
