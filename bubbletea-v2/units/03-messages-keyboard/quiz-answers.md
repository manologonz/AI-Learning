# Unit 3 Quiz — Answers

**Q1: B.** `KeyPressMsg` is the press event. A (interface) also matches releases; C doesn't exist; D is the data struct, not a message.

**Q2: B.** `"space"`. A is v1 behavior. C is invented. D would be a bug-hiding empty string, but not what `String()` returns.

**Q3: B.** `Text` is `""` for non-typing keys and the typed string otherwise. A would append `"up"`, `"enter"` etc. C converts arrow constants to garbage runes. D includes modifier notation (`"ctrl+a"`).

**Q4: B.** Initially and on each resize. A misses the first size; C is not how it works (`RequestWindowSize` exists but isn't required); D is nonsense.

**Q5: C.** It's an interface implemented by both `KeyPressMsg` and `KeyReleaseMsg`; mouse events are separate types.

**Q6:** `msg.Type` → `msg.Code` (a `rune`), `msg.Runes` → `msg.Text` (a `string`).

**Q7:** Legacy terminals encode both as byte 0x09 (ctrl+letter maps to control codes), so the app can't tell them apart unless a keyboard-enhancement protocol is active.

**Q8:**
```go
type totalMsg struct{ n int }
// in Update's type switch:
case totalMsg:
    m.total = msg.n
```

**Q9:** `WindowSizeMsg` hasn't arrived yet, so `m.w` is 0. Fix: guard: `if m.w == 0 { return tea.NewView("") }` or default sizes in the constructor.

**Q10: A.** Raw mode delivers ctrl+c as a key; you choose to quit. B isn't reliable in raw mode; C skips terminal restoration; D `tea.Interrupt` is a different mechanism that aborts the program with `ErrInterrupted`, not the idiomatic quit.
