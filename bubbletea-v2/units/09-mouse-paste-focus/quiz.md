# Unit 9 Quiz

**Q1 (easy).** Mouse events are:
A) On by default  B) Enabled via `v.MouseMode`  C) Enabled via `tea.WithMouse()`  D) Only in alt-screen

**Q2 (easy).** Which message type reports a wheel scroll?
A) `MouseClickMsg`  B) `MouseWheelMsg`  C) `MouseScrollMsg`  D) `MouseMotionMsg`

**Q3 (medium).** `msg.Y` in a `MouseClickMsg` is relative to:
A) The top of your view  B) The top-left of the terminal  C) The bottom of the screen  D) The list widget

**Q4 (medium).** With bracketed paste, pasted text arrives as:
A) Individual `KeyPressMsg`s  B) One `PasteMsg`  C) A `ClipboardMsg`  D) An `EnvMsg`

**Q5 (medium).** To be told when the terminal loses focus:
A) `v.ReportFocus = true` then handle `BlurMsg`  B) It always arrives  C) `tea.WithFocus()`  D) Poll `os.Stdin`

**Q6 (medium).** Short answer: Why should a mouse UI also work by keyboard?

**Q7 (hard).** Short answer: your handler on `tea.KeyMsg` fires twice per keystroke after you set `ReportEventTypes = true`. Why and fix?

**Q8 (medium).** Which mouse mode reports hover motion?
A) None  B) CellMotion  C) AllMotion  D) Hover

**Q9 (hard, apply).** Write the `Update` case mapping a left-click to a list row index given `topRow` and scroll `offset`.

**Q10 (medium).** Why is `tea.SetClipboard` unreliable as the only "copy" mechanism?
