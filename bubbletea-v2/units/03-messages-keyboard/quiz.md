# Unit 3 Quiz

**Q1 (easy).** In v2 which type do you normally match for "the user pressed a key"?
A) `tea.KeyMsg`  B) `tea.KeyPressMsg`  C) `tea.KeyEvent`  D) `tea.Key`

**Q2 (easy).** `msg.String()` for the space bar?
A) `" "`  B) `"space"`  C) `"spc"`  D) `""`

**Q3 (medium).** Which expression best appends typed characters to a buffer without inserting junk for arrow keys?
A) `m.text += msg.String()`  B) `m.text += msg.Key().Text`  C) `m.text += string(msg.Key().Code)`  D) `m.text += msg.Keystroke()`

**Q4 (medium).** When is `WindowSizeMsg` sent?
A) Only on resize  B) Once initially and on every resize  C) Only if requested in `Init`  D) Every frame

**Q5 (medium).** What does `case tea.KeyMsg:` match in v2?
A) Only presses  B) Only releases  C) Both presses and releases  D) Mouse events too

**Q6 (medium).** Short answer: name the v2 replacements for v1's `msg.Type` and `msg.Runes`.

**Q7 (hard).** Short answer: why are `ctrl+i` and `tab` usually indistinguishable in terminals?

**Q8 (medium, apply).** Write a custom message type and the `case` that stores its payload (`n int`) into `m.total`.

**Q9 (hard).** In `View` you compute `m.w / 2` and get a divide/zero-width layout on the first frame. Why, and what's the fix?

**Q10 (medium).** Which case is the correct way to quit on Ctrl+C?
A) `case tea.KeyPressMsg: if msg.String()=="ctrl+c" {return m, tea.Quit}`  B) Rely on SIGINT  C) `os.Exit(1)`  D) `return m, tea.Interrupt` only
