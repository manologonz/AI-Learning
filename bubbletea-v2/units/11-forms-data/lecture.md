# Unit 11: Forms, Validation & Data-Heavy UIs

> Runnable code: `examples/u11-form/main.go` (hand-built form + table) and `examples/u11-huh/main.go` (the `huh` library embedded in a model).

## Objectives
By the end of this unit you can:
- Build a multi-field form with focus movement, per-field validation, and submit logic.
- Decide when to hand-roll a form vs use `huh`, and embed a `huh.Form` in your own model.
- Present large datasets with `table`, `list`, and `viewport` without rendering everything.
- Implement filtering, sorting, and selection state, and keep the selection stable across data refreshes.
- Represent loading / empty / error / loaded states explicitly (a state machine, not booleans).

## Prerequisites
Units 7, 8, 4.

## Core content

### Hand-built forms: the structure
```go
type model struct {
    inputs []textinput.Model   // 0=name, 1=email
    focus  int
    err    error
}
```
Ingredients (all from Unit 7):
- **Focus movement**: `tab`/`down` → next, `shift+tab`/`up` → previous. `Blur()` the old input; `Focus()` the new and **return its Cmd**.
- **Forward keys only to the focused input** (`m.inputs[m.focus], cmd = m.inputs[m.focus].Update(msg)`).
- **Validation**: `textinput.Validate func(string) error` runs as the value changes; the error is in `input.Err`. Show it under the field in an error style.
- **Submit**: on `enter` on the last field (or a dedicated key), validate *everything* (don't trust per-field state alone: an untouched field never ran its validator), then emit a `submitMsg` or update state.

See `examples/u11-form/main.go`: two fields, `mail.ParseAddress` validation, and a table showing submitted rows.

### Validation principles
1. **Validate at the right time**: live validation while typing is annoying for "email" (it's invalid until the last char). Prefer validating on blur/submit for format checks; live-validate only for character limits and allowed chars.
2. **Separate checking from presenting**: validator returns `error`; the view decides styling and placement.
3. **Server-side truth**: uniqueness ("username taken") needs a command; represent it as `checking → ok/taken` field state and ignore stale replies (Unit 4 IDs).
4. **Empty is a state**: treat "" separately from "invalid" (required vs format).

### Using `huh`
`charm.land/huh/v2` is a form library built on Bubble Tea: groups of fields (`Input`, `Text`, `Select`, `MultiSelect`, `Confirm`, …), themes, validation, accessibility mode.
- **Standalone**: `err := huh.NewForm(groups...).Run()` blocks and fills the variables you bound with `.Value(&x)`. Perfect for CLIs that just need to ask questions.
- **Embedded in your model** (see `examples/u11-huh`): a `*huh.Form` is itself a `tea.Model`-like component:
  ```go
  form, cmd := m.form.Update(msg)
  if f, ok := form.(*huh.Form); ok { m.form = f }     // it returns tea.Model: assert back
  if m.form.State == huh.StateCompleted { /* read the bound vars */ }
  // Init: return m.form.Init();  View: m.form.View()
  ```
Choose `huh` when the UI is "ask several questions in sequence"; hand-roll when the form is part of a bespoke layout, needs unusual interactions, or you need total control of keys and rendering. Note the type assertion: in Update, the form's returned `tea.Model` is the same pointer, but assert it rather than assuming.

### Big data: don't render what you can't see
`View` runs after every message. Building a 100,000-row string every keypress will lag.
- **Window it**: render only `height` rows starting at `offset`. Bubbles `list`, `table`, and `viewport` already do this (`table.WithHeight`, `list.SetSize`).
- **Paginate** long lists (`list` paginates; `paginator` is the standalone bubble).
- **Precompute** expensive per-row strings (formatted columns) when data changes, not on every `View`.
- **Load asynchronously** (Unit 4): fetch pages via commands and append.

### Filtering and sorting
- `list` has built-in fuzzy filter (`/`). For your own datasets: keep `all []row` (source of truth) and `visible []row` (derived: filtered+sorted). Recompute `visible` *in `Update`* when the filter/sort changes (or lazily with a `dirty` flag), never mutate `all`.
- **Stable selection across refreshes**: store the selected *ID*, not the *index*. After the data changes, find the ID in the new `visible` slice and re-select it; fall back to the nearest index if it's gone.

### Explicit data states
Booleans (`loading`, `err != nil`, `len(items)==0`) invite impossible combinations ("loading and error"). Prefer one field:
```go
type state int
const ( stateLoading state = iota; stateError; stateEmpty; stateReady )
```
`View` switches on `state`; every state has a designed screen (spinner / error + retry key / "nothing here yet" / the table). Transitions happen only in `Update` in response to messages.

### Editing data in tables/lists
Bubbles' `table` displays `[]string` rows; keep your own `[]Record` and derive rows: `rows[i] = table.Row{r.Name, r.Email}`. Editing = update the record, then call `t.SetRows(rows)`. Do not treat table rows as your database.

## Worked examples

### Straightforward: one field with a validator
```go
in := textinput.New()
in.Validate = func(s string) error {
    if len(s) < 3 { return errors.New("min 3 chars") }
    return nil
}
// View: in.View() + (if in.Err != nil: render error)
```

### Typical: a filterable table with stable selection
Model: `all []Record`, `filter textinput.Model`, `table table.Model`, `selectedID string`. On filter change: `visible := filter(all, query)`; `rows := toRows(visible)`; `t.SetRows(rows)`; restore cursor with `t.SetCursor(indexOfID(visible, selectedID))` (falling back to 0). On cursor move, read `t.SelectedRow()` (or index into `visible`) and update `selectedID`.

### Edge case: validating an untouched form, and pasting into a single-line field
Untouched inputs never triggered `Validate`, so `in.Err == nil` even though the value is empty: submit logic must check `Value() == ""` explicitly (as `examples/u11-form` does). Second: pasting a multi-line string into a single-line input: `textinput` sanitizes newlines, but your *own* field types must too (Unit 9). Third: on very long tables, `SetRows` on every keystroke of a filter re-renders everything; debounce or filter in a command for huge data.

## Common misconceptions
- **"`Err == nil` means the field is valid."** Only if the validator ran; check emptiness on submit.
- **"Store the selected index."** Indexes shift when data changes; store an ID.
- **"Render everything and let the terminal scroll."** It makes every keypress slow.
- **"A table row is my data model."** Keep records separately.
- **"`huh` replaces Bubble Tea."** It's built on it; embed it as a component.
- **"loading/error/empty can be independent booleans."** Model states explicitly.

## Check your understanding
1. Why validate everything on submit, not just rely on `Err`? *(Untouched fields never validated.)*
2. What do you store to keep selection across refreshes? *(The record ID.)*
3. When is `huh.NewForm(...).Run()` the right tool? *(A CLI asking a few questions.)*
4. Why window large lists? *(View cost per message.)*

## Summary
- Forms = focus management + forwarding to the focused field + validation + submit checks.
- `huh` for sequential questions; hand-roll for custom layouts.
- Render only visible rows; keep `all` vs `visible`; select by ID.
- Model loading/error/empty/ready as one explicit state.
