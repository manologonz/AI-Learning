# Unit 7 Quiz

**Q1 (easy).** Bubbles v2 import path prefix?
A) `github.com/charmbracelet/bubbles`  B) `charm.land/bubbles/v2`  C) `charm.land/bubbles`  D) `bubbles.charm.land/v2`

**Q2 (easy).** What does a component's `View()` return?
A) `tea.View`  B) `string`  C) `[]string`  D) `tea.Model`

**Q3 (medium).** Which line correctly forwards to a text input?
A) `m.input.Update(msg)`  B) `m.input, cmd = m.input.Update(msg)`  C) `cmd = m.input.Update(msg)`  D) `m.input = m.input.Update`

**Q4 (medium).** What happens if you drop the `Cmd` returned by `textinput.Update`/`Focus()`?
A) Nothing  B) The cursor blink (and other async behavior) stops  C) Panic  D) Input is lost

**Q5 (medium).** In v2, set a textinput's width with:
A) `ti.Width = 30`  B) `ti.SetWidth(30)`  C) `textinput.WithWidth(30)` only  D) `ti.Size(30)`

**Q6 (medium).** Short answer: why keep forwarding `spinner.TickMsg` even when the spinner's area isn't focused?

**Q7 (hard).** Short answer: describe the bug when a global `q`→quit binding is checked before forwarding to a `list` in filter mode, and the fix.

**Q8 (medium).** `key.Binding` benefits — pick the false statement:
A) One place to change keys  B) Powers `help` text  C) Can be disabled with `SetEnabled(false)`  D) Makes keys work in every terminal

**Q9 (hard, apply).** Write the `list.Item` implementation for `type task struct{ name string; done bool }` and the code retrieving the selected task.

**Q10 (medium).** Why do v2 style helpers like `help.DefaultStyles` take `isDark`?
