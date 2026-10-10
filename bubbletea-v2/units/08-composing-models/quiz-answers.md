# Unit 8 Quiz — Answers

**Q1: B.** Messages keep children decoupled. A couples; C is hidden shared state (and racy); D kills the app.

**Q2: B.** Global concerns first, then routing, then delegate. A duplicates code; C/D can't handle input.

**Q3: B.** Key presses go only to the focused child. Non-key messages drive animations and async completion for any child.

**Q4: B.** You need assertions (which can panic) and lose concrete types. They *can* render, receive messages, return commands.

**Q5:** The parent: it lays out all children, so it can allocate rectangles; children stay reusable and testable at any size.

**Q6:** The form model isn't reinitialized on entry (`m.form = newForm()`), so old state persists.

**Q7:** (1) Deliver the result to A's model regardless of which screen is active (store it; show it when the user returns): route by message *type*. (2) Discard stale/irrelevant results: include an ID and ignore mismatches. Implementation: `case resultMsg: if msg.id == m.pendingID { m.a.data = msg.data }`.

**Q8:**
```go
type navMsg struct{ to screen }
func goSettings() tea.Msg { return navMsg{settingsScreen} }
// in child: return c, goSettings     (a Cmd is func() Msg)
```

**Q9:** The returned child model is discarded, so nothing changes; also the child's command is dropped. Must be `m.form, cmd = m.form.Update(msg)`.

**Q10:** Two sources of truth can diverge: updating one leaves the other stale, causing inconsistent UI. Keep each fact in exactly one model.
