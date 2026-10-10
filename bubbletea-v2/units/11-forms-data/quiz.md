# Unit 11 Quiz

**Q1 (easy).** Where does `textinput.Validate`'s error appear?
A) In `input.Err`  B) It panics  C) In the next `Update` as `errMsg`  D) In stdout

**Q2 (easy).** To keep a selection across data refreshes you should store the selected:
A) Index  B) ID  C) Screen row  D) Pointer to the table

**Q3 (medium).** Which is the best default representation of screen data states?
A) Three booleans  B) A single `state` enum  C) Strings  D) A channel

**Q4 (medium).** The standalone `huh.NewForm(...).Run()`:
A) Needs your own `tea.Program`  B) Runs its own program and blocks until done  C) Returns immediately  D) Only works over SSH

**Q5 (medium).** Why can an untouched required field show `Err == nil`?
A) Validators only run when the value changes  B) Validators are broken  C) Empty always valid  D) `Err` is only for pastes

**Q6 (medium).** Short answer: give two reasons not to render all 100k rows in `View`.

**Q7 (hard).** Short answer: describe the fields you'd keep for a filterable list (`all`, `visible`, …) and when to recompute.

**Q8 (medium).** True/False: table rows should be your source-of-truth data.

**Q9 (hard, apply).** Write a validator func for "username: 3–12 chars, letters/digits only".

**Q10 (medium).** Why should uniqueness checks use a command and IDs on replies?
