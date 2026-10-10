# Unit 4 Exam

**E1 (Apply).** Write a spinner-less "loading…" flow: on `r`, set `loading=true` and fetch a string via `fetch(url)`; when done show it; on error show the error. Give the model, `Update`, and `fetch`.

**E2 (Analyze).** This clock is too fast after pressing `s`. Find the cause and fix it.
```go
func (m model) Init() tea.Cmd { return tick() }
case tea.KeyPressMsg: if msg.String()=="s" { return m, tick() }
case tickMsg: m.n++; return m, tick()
```

**E3 (Evaluate).** Choose `Batch` or `Sequence` and justify: (a) fetch users and fetch settings on start; (b) write config, then print "saved", then quit; (c) start two spinners.

**E4 (Analyze).** Explain the data race in:
```go
return m, func() tea.Msg { m.cache["k"] = compute(); return doneMsg{} }
```
with `cache` a map in the model. Give the fix.

**E5 (Apply).** Add stale-reply protection: a `searchMsg{seq int, results []string}`, a model field `seq`, and code for issuing and handling.

**E6 (Understand).** Draw (in text) the message flow for a keypress that triggers an HTTP call and shows the result. Label who runs each step and on which goroutine.

**E7 (Apply).** Implement a 10-second countdown: fields, `Init`, tick handling, stop at 0, and show "Done!" then quit after 1 more second.

**E8 (Analyze).** Why is `tea.Batch(nil, nil)` harmless, and what does `Batch(cmd)` with one command do?

**E9 (Evaluate).** A teammate wraps every state change in a Cmd "for purity", e.g. `return m, func() tea.Msg { return incMsg{} }` for `count++`. Critique.

**E10 (Create).** Design the messages and flow for auto-save: save 2 seconds after the *last* keystroke (debounce). Describe the messages, model fields, and how stale timers are ignored.
