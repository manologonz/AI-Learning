# Unit 12 Quiz — Answers

**Q1: B.** The `/v2` module matches Bubble Tea v2. A is the v1 path; C/D don't exist.

**Q2: B.** stdout/stderr/View collide with the renderer.

**Q3: B.** Call the `Cmd` and check for `tea.QuitMsg`. A the model doesn't record quitting; C nothing is printed; D false.

**Q4: B.** Condition + timeout; sleeps are slow and flaky.

**Q5: A.** `FinalModel` waits for the program to finish; if nothing sends `q`/quit, it times out (with `WithFinalTimeout` it fails instead of hanging forever).

**Q6:** Inject the fetch function (or an `http.Client`/`httptest.Server`) into the model or command builder; call the `Cmd` and assert the returned `Msg` for both success and error paths.

**Q7:** `msg` is a `tea.BatchMsg` (`[]tea.Cmd`). Type-assert it, then call each cmd: `for _, c := range msg.(tea.BatchMsg) { r := c(); … }`.

**Q8: B.** `TEA_DEBUG=true` writes `bubbletea-panic-<unix>.log`.

**Q9:**
```go
tests := []struct{ name string; in int; key string; want int }{
    {"x", 1, "x", 2}, {"z", 5, "z", 0},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        next, _ := model{count: tt.in}.Update(tea.KeyPressMsg{Code: rune(tt.key[0]), Text: tt.key})
        if got := next.(model).count; got != tt.want { t.Fatalf("got %d", got) }
    })
}
```

**Q10:** Animations (spinner frames, timestamps), color-profile/terminal-size differences, map iteration order, random data.
