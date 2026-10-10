# Unit 4 Exam — Answers

**E1:**
```go
type model struct{ loading bool; body string; err error }
type bodyMsg string
type errMsg struct{ err error }
func fetch(url string) tea.Cmd {
    return func() tea.Msg {
        res, err := http.Get(url)            // better: use a context with timeout
        if err != nil { return errMsg{err} }
        defer res.Body.Close()
        b, err := io.ReadAll(res.Body)
        if err != nil { return errMsg{err} }
        return bodyMsg(b)
    }
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyPressMsg:
        if msg.String() == "r" { m.loading, m.err = true, nil; return m, fetch("https://example.com") }
    case bodyMsg: m.loading, m.body = false, string(msg)
    case errMsg:  m.loading, m.err = false, msg.err
    }
    return m, nil
}
```
Points: error path, resetting `loading`, closure over `url`.

**E2:** Pressing `s` starts a second tick chain; each chain re-arms itself forever, so `n` increments twice per second. Fix: only start a tick if not running (`if !m.running { m.running = true; return m, tick() }`) or remove the `s` handler.

**E3:** (a) Batch: independent, order doesn't matter. (b) Sequence: order matters (write → message → quit). (c) Batch: independent (spinners are usually just one `Tick` each and don't care about order).

**E4:** The command runs on another goroutine while `Update`/`View` read the same map (maps share memory across model copies): concurrent map write/read → race/panic. Fix: compute in the command, return the result in a message (`computedMsg{k, v}`), and let `Update` write to the map.

**E5:**
```go
type searchMsg struct{ seq int; results []string }
func search(q string, seq int) tea.Cmd {
    return func() tea.Msg { return searchMsg{seq, doSearch(q)} }
}
// on keystroke:
m.seq++; return m, search(m.query, m.seq)
// on reply:
case searchMsg:
    if msg.seq != m.seq { return m, nil } // stale
    m.results = msg.results
```
Note: `m.seq` is incremented *before* the closure is built and `seq` is passed by value.

**E6:** `KeyPressMsg` → (event loop goroutine) `Update` sets loading, returns `fetchCmd` → (runtime goroutine) executes HTTP → returns `resultMsg` → (event loop) `Update` stores result → `View` re-renders. Key idea: only the HTTP step happens off-loop.

**E7:**
```go
type model struct{ left int }
func (m model) Init() tea.Cmd { return tick() }
case tickMsg:
    if m.left > 0 { m.left--; if m.left == 0 { return m, tea.Tick(time.Second, func(time.Time) tea.Msg { return quitMsg{} }) }; return m, tick() }
case quitMsg: return m, tea.Quit
```
with `model{left: 10}`; `View` shows `"Done!"` when `left == 0`. Alternative: `tea.Sequence(tea.Tick(...), tea.Quit)`.

**E8:** `Batch` discards nil commands, so `Batch(nil,nil)` yields `nil`; `Batch(cmd)` with a single command returns effectively just that command (no extra overhead). Harmless because "no work" is representable.

**E9:** Needless indirection: it makes each increment a round-trip through the runtime (extra render/latency, harder to test and reason about); commands are for I/O and time, not internal state changes. Do `m.count++` directly (Bubble Tea docs: "there's almost never a reason to use a command to send a message to another part of your program").

**E10:** Fields: `edits int`. On each keystroke: `m.edits++; id := m.edits; return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg { return saveDueMsg{id} })`. On `saveDueMsg`: `if msg.id != m.edits { return m, nil }` (a newer keystroke superseded it) else return `saveCmd(text)`. On `savedMsg`/`errMsg`: update status. This is the sequence-number pattern from the lecture applied to timers.
