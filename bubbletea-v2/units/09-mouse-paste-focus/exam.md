# Unit 9 Exam

**E1 (Apply).** Implement click-to-select for a two-pane layout: left pane width 30, right pane the rest; a click sets `m.focus` to `left` or `right`. Include the View settings needed.

**E2 (Analyze).** Your inline (non-alt) mouse UI selects the wrong row by exactly 12 lines. Why? Give two fixes.

**E3 (Evaluate).** For a drawing app, pick `CellMotion` or `AllMotion` and justify; for a hover tooltip UI, which and why?

**E4 (Analyze).** A search box with `q` = quit shortcut quits when the user pastes text containing `q`. Diagnose in both a bracketed-paste-capable and non-capable terminal, and fix what you can.

**E5 (Apply).** Add pause-on-blur to a ticking clock: fields, messages, and how the tick chain is stopped and restarted without duplicating chains.

**E6 (Evaluate).** Design a fallback strategy for a "hold Space to talk" feature that needs key release events.

**E7 (Analyze).** Explain why `tea.MouseMsg` (interface) is useful and give a case where you'd choose it over the four concrete types.

**E8 (Apply).** Implement a draggable divider: on click near the divider start dragging, on motion update width, on release stop. State the mouse mode and messages.

**E9 (Evaluate).** Argue for or against enabling `AllMotion` by default in a list UI.

**E10 (Create).** Specify the input handling design (messages, modes, fallbacks) for a terminal file manager with clickable rows, wheel scroll, paste into a path field, and copy-path-to-clipboard.
