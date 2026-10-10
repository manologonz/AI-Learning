# Unit 11 Quiz — Answers

**Q1: A.** The field is `input.Err`. Not a panic, message or stdout.

**Q2: B.** IDs are stable; A/C shift with filtering/sorting; D isn't meaningful.

**Q3: B.** One enum prevents impossible combinations. A allows "loading and error"; C loses type-safety; D is irrelevant.

**Q4: B.** It creates and runs its own program. A is for the embedded style; C/D are false.

**Q5: A.** Validation is triggered by changes; an untouched empty input never ran it, so check emptiness at submit.

**Q6:** `View` runs on every message so cost multiplies per keystroke (lag); huge strings also cause heavy diffing/allocation and terminal output; the terminal only shows a screenful anyway.

**Q7:** `all []Record` (source of truth), `visible []Record` (filtered/sorted, derived), `query`, `sortKey`, `selectedID`. Recompute `visible` in `Update` when `query`/`sortKey`/`all` changes (not in `View`).

**Q8:** False. Rows are a view of your data; keep records separately.

**Q9:**
```go
func validateUser(s string) error {
    if l := len(s); l < 3 || l > 12 { return errors.New("3–12 characters") }
    for _, r := range s { if !unicode.IsLetter(r) && !unicode.IsDigit(r) { return errors.New("letters and digits only") } }
    return nil
}
```
(Checks `len` on bytes; use `utf8.RuneCountInString` if non-ASCII counts matter.)

**Q10:** It's I/O (must not block `Update`); replies may arrive out of order or after the value changed, so tagging replies with an ID/value lets you ignore stale ones.
