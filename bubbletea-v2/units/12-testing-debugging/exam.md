# Unit 12 Exam

**E1 (Apply).** Write tests for the u12 counter: `up` increments, `down` doesn't go below 0, `q` sets `done` and returns a quit command.

**E2 (Apply).** Refactor a model with a hard-coded `http.Get` in a command so it's testable; write two tests (success + failure).

**E3 (Analyze).** This teatest test hangs until the `go test` timeout. Why? Fix it.
```go
tm := teatest.NewTestModel(t, model{})
tm.Send(tea.KeyPressMsg{Code: tea.KeyUp})
fm := tm.FinalModel(t).(model)
```

**E4 (Analyze).** A golden file test passes locally and fails in CI. List three plausible causes and fixes.

**E5 (Evaluate).** Which tests would you write for the Unit 8 router (menu/form/done)? Give the pyramid: which level for each behavior.

**E6 (Apply).** Write a `WaitFor` condition that waits for the text "Loaded 3 items" and fails after 3s.

**E7 (Analyze).** Debugging: pressing `Space` doesn't toggle a checkbox. You log `%T %+v` for messages and see `tea.KeyPressMsg{Code:32 Text:" "}`. Where's the bug likely?

**E8 (Evaluate).** When is a teatest end-to-end test worth its cost over Update tests?

**E9 (Analyze).** Explain how you'd test that a stale reply (Unit 4's sequence numbers) is ignored.

**E10 (Create).** Design the test suite (files, helpers, levels) for the capstone app so that CI runs in <5s without a TTY.
