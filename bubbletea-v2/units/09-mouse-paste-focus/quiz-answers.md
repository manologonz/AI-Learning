# Unit 9 Quiz — Answers

**Q1: B.** A is false; C was removed with v1 options; D is false (works inline too, with caveats).

**Q2: B.** `MouseWheelMsg`; `MouseScrollMsg` doesn't exist.

**Q3: B.** Zero-based from the terminal's top-left. In inline mode it isn't your view's origin.

**Q4: B.** One `PasteMsg` (with `PasteStartMsg`/`PasteEndMsg` around it).

**Q5: A.** Opt in; then `FocusMsg`/`BlurMsg` arrive. B is false; C/D don't exist.

**Q6:** Accessibility, SSH/tmux/terminal limitations, muscle memory, and some terminals swallow mouse events; keyboard parity is the baseline.

**Q7:** `tea.KeyMsg` is an interface matching both press and release; with release events enabled each keystroke generates two messages. Match `tea.KeyPressMsg` instead.

**Q8: C.** `CellMotion` reports clicks/drags/wheel only.

**Q9:**
```go
case tea.MouseClickMsg:
    if msg.Button == tea.MouseLeft {
        if idx := msg.Y - topRow + m.offset; idx >= 0 && idx < len(m.items) { m.cursor = idx }
    }
```

**Q10:** OSC 52 support and policies vary (some terminals disable it, block reads, or limit size), and there's no delivery confirmation; offer an alternative (e.g. write to a file or show text to copy).
