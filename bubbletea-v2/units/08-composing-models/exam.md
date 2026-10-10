# Unit 8 Exam

**E1 (Apply).** Extend `examples/u08-screens` with a "confirm quit" screen: choosing Quit shows "Really quit? y/n"; `y` quits, `n` returns to the menu. Describe the new screen constant, child model, messages, and router changes.

**E2 (Analyze).** Two `textinput`s both show typed characters simultaneously. Give two possible causes and fixes.

**E3 (Analyze).** After adding a spinner in a child that's inactive half the time, the spinner freezes when you return. Diagnose.

**E4 (Evaluate).** Compare *concrete-typed children* to *`tea.Model` children* for a wizard with 5 steps of different types. Which would you choose and why?

**E5 (Apply).** Implement window-resize distribution for a three-pane UI: sidebar (fixed 24), main (fill), footer (1 line).

**E6 (Analyze).** Explain why a child holding `parent *model` is problematic in a value-receiver Bubble Tea app.

**E7 (Apply).** Write focus management for `[]tea.Model`-like children with `tab`/`shift+tab`, and show which messages get forwarded to whom.

**E8 (Evaluate).** A teammate puts *all* app state (screen data, focus, forms) in one giant struct with a single Update. List costs and benefits; propose a refactor.

**E9 (Analyze).** Find the bug: 
```go
case navMsg: m.current = msg.to
case tea.KeyPressMsg: m.form, cmd = m.form.Update(msg); m.menu, cmd = m.menu.Update(msg)
return m, cmd
```

**E10 (Create).** Design a "modal dialog" overlay on top of any screen: how it captures input, how it returns a result, and how it prevents the underlying screen from reacting.
