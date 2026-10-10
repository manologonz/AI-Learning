# Unit 4: Commands — Side Effects Without Breaking Purity

> Runnable code: `examples/u04-commands/main.go` (ticking clock + HTTP check).

## Objectives
By the end of this unit you can:
- Explain what a `tea.Cmd` is and why `Update` returns one instead of doing I/O.
- Write commands that do work and return messages (including error messages).
- Use `tea.Tick`/`tea.Every` to build repeating timers (re-arming from `Update`).
- Choose correctly between `tea.Batch` (concurrent, unordered) and `tea.Sequence` (ordered).
- Avoid the classic bugs: blocking in `Update`, sharing mutable state with commands, forgetting to re-arm a tick.

## Prerequisites
Units 1–3 (closures, messages).

## Core content

### The purity problem
`Update` and `View` should be fast and deterministic. But real apps must fetch URLs, read files, wait for timers. **Blocking inside `Update` freezes the whole UI**: the event loop is single-threaded (Unit 2).

The Elm answer: `Update` doesn't *perform* effects; it *describes* them by returning a **command**. The runtime performs the command on its own goroutine and feeds the result back as a **message**.

```go
type Cmd func() Msg     // "a function that, when run, eventually produces a message"
```
(recall Unit 1: closures make this natural). Returning `nil` means "no effect".

```
Update ──returns──► Cmd ──runtime runs it in a goroutine──► Msg ──► Update …
```

### Writing a command
A command is usually a function that *returns* a `Cmd`, so you can pass arguments (via closure):
```go
type statusMsg struct{ code int }
type errMsg struct{ err error }
func (e errMsg) Error() string { return e.err.Error() }

func checkURL(url string) tea.Cmd {
    return func() tea.Msg {                    // runs OFF the event loop
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
        if err != nil { return errMsg{err} }
        res, err := http.DefaultClient.Do(req)
        if err != nil { return errMsg{err} }   // errors are just messages
        defer res.Body.Close()
        return statusMsg{res.StatusCode}
    }
}
```
Usage: `case "enter": m.loading = true; return m, checkURL("https://example.com")`.
Then handle `statusMsg`/`errMsg` in `Update` to clear `loading` and store the result.

**Golden rules**
1. The command function must not touch the model. It runs concurrently with the loop. Pass it *values* it needs (copy them into the closure) and return results as messages.
2. Always have an error path: `errMsg`. A command that can fail but can't say so leaves the UI stuck on "loading…".
3. Always set a timeout/context on network calls.

### Built-in commands
- `tea.Quit`: exit gracefully.
- `tea.Batch(cmds...)`: run commands **concurrently**; messages arrive in whatever order they finish.
- `tea.Sequence(cmds...)`: run **one after another**, in order (v1's `Sequentially` was renamed).
- `tea.Tick(d, fn)`: after duration `d`, run `fn(time)` to make a message. **Fires once.**
- `tea.Every(d, fn)`: like `Tick` but aligned to the system clock (e.g. top of the second).
- `tea.Println/Printf`: print a line *above* the TUI in non-alt-screen mode (persistent log lines).
- `tea.RequestWindowSize`, `tea.Exec/ExecProcess`, `tea.Suspend` (Unit 13).

### Ticking: the re-arm pattern
`Tick` sends **one** message. For a clock, return a new tick each time you handle one:
```go
type tickMsg time.Time
func tick() tea.Cmd {
    return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}
func (m model) Init() tea.Cmd { return tick() }          // start
// in Update:
case tickMsg:
    m.secs++
    return m, tick()                                     // re-arm; omit to stop
```
Stopping is free: just don't return another tick (e.g. when `paused`).

### Batch vs Sequence
```go
return m, tea.Batch(fetchUsers(), fetchPosts())          // parallel, arrival order undefined
return m, tea.Sequence(saveFile(), tea.Quit)             // save FIRST, then quit
```
Use `Sequence` whenever step B must not begin until A's command has finished (save-then-quit is the canonical case; with `Batch` the quit could race the save). Multiple nil cmds inside `Batch` are fine; `Batch` with a single non-nil returns it unchanged.

### Correlating replies
Concurrent commands can return in any order, and stale replies can arrive after the user moved on. Give messages **identity**: an ID/sequence number in the message, compare in `Update`:
```go
type resultMsg struct{ id int; body string }
case resultMsg:
    if msg.id != m.latestID { return m, nil }   // ignore stale reply
```

## Worked examples

### Straightforward: a one-shot timer
```go
func delayedHello() tea.Cmd {
    return tea.Tick(2*time.Second, func(time.Time) tea.Msg { return helloMsg{} })
}
```
Press a key → return `delayedHello()` → 2s later `helloMsg` arrives → `Update` sets `m.greeting`.

### Typical: load data on startup, with a loading flag and an error state
```go
type model struct { loading bool; items []string; err error }
func (m model) Init() tea.Cmd { return loadItems }          // Init returns the Cmd itself
func loadItems() tea.Msg {                                   // a Cmd is just func() Msg
    items, err := readFromDisk()
    if err != nil { return errMsg{err} }
    return itemsMsg(items)
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case itemsMsg: m.items, m.loading = msg, false
    case errMsg:   m.err, m.loading = msg.err, false
    }
    return m, nil
}
```
Note `loadItems` has the signature `func() tea.Msg`, so **the function itself is a Cmd**; no wrapper required when it needs no arguments.

### Edge case: commands run later, on another goroutine
```go
// RISKY when the model holds reference types (maps, slices, pointers):
return m, func() tea.Msg { return search(m.query, m.cache) }  // reads m.cache on the command goroutine
```
The closure runs *later*, concurrently with the event loop. Value fields like `m.query` were already copied when the model was copied, but a map like `m.cache` is shared, so `Update` writing it while the command reads it is a data race. **Rule: copy the plain values you need into locals before building the closure, and never share mutable reference fields with a command.** (`q := m.query; return m, func() tea.Msg { return search(q) }`; results come back as messages.)

Second edge: the lost re-arm. A `tickMsg` case that *forgets* `return m, tick()` silently stops the clock after one second. And the reverse: returning `tick()` from *both* the `tickMsg` handler and elsewhere (e.g. each keypress) creates **multiple ticking chains**, so the clock runs twice as fast. Guard with a `running` flag.

## Common misconceptions
- **"A Cmd runs when I return it *from* `Update` immediately and synchronously."** It runs asynchronously on a goroutine; `Update` returns first.
- **"I can mutate the model from inside the command."** No; return a message and let `Update` do it.
- **"`Tick` repeats."** It fires once; you re-arm.
- **"`Batch` runs in order."** It's concurrent; use `Sequence` for order.
- **"Errors from commands crash the app."** They vanish unless you turn them into messages.
- **"I should sleep in `Update` to wait."** That freezes rendering and input.

## Check your understanding
1. What's the type of `tea.Cmd`? *(`func() tea.Msg`)*
2. You want "save then quit". `Batch` or `Sequence`? *(Sequence.)*
3. Your clock runs at double speed after pressing a key that also returns `tick()`. Why? *(Two tick chains.)*
4. Where should the HTTP call happen? *(Inside a Cmd, never in `Update`.)*

## Summary
- Commands are how `Update` requests side effects; results return as messages.
- Commands run concurrently: don't touch the model; capture values; always return `errMsg` on failure.
- `Tick` is one-shot; re-arm in `Update`. `Batch` = parallel, `Sequence` = ordered.
- Give messages IDs to discard stale replies.
