# Unit 14 Exam — Answers

**E1:** `cmd/notes/main.go` (flags/wiring; imports ui, core, ssh server), `internal/core` (Note, Store interface, search; imports stdlib only), `internal/storage` (file/SQLite implementing `core.Store`), `internal/ui` (Bubble Tea models; imports core, bubbletea/bubbles/lipgloss), `internal/server` (Wish; imports ui). `list --json` imports core+storage only.

**E2:**
```go
func main() {
    jsonOut := flag.Bool("json", false, "print JSON and exit"); ver := flag.Bool("version", false, "")
    flag.Parse()
    if *ver { fmt.Println(version); return }
    interactive := term.IsTerminal(os.Stdin.Fd()) && term.IsTerminal(os.Stdout.Fd())
    if *jsonOut || !interactive { if err := printPlain(*jsonOut); err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }; return }
    if _, err := tea.NewProgram(ui.New()).Run(); err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
}
```
Bad flags → `flag` exits 2 automatically.

**E3:** A model (or a pointer inside it) is shared across sessions (global model/slice/map). Create the model inside the Wish handler per session; share only via synchronized external state or `Send`.

**E4:** `flag` if just 2–3 flags; Cobra if you have subcommands (`list`, `serve`, `export`), completions, and help generation. For this tool with `list`/`serve`/TUI default, Cobra is reasonable.

**E5:** `wish.WithMiddleware(wishtea.Middleware(handler), activeterm.Middleware(), logging.Middleware())`. Wish composes middleware first-to-last, so **the last listed runs first**: logging sees every connection, then `activeterm` rejects sessions without a PTY, then the Bubble Tea handler runs. Then `ListenAndServe` in a goroutine, `signal.Notify` for SIGINT/SIGTERM, and `Shutdown(ctx)` with a timeout.

**E6:** stdout is a pipe: the TUI writes escape sequences to it and raw mode may fail/block; `head` closes early. Fix: detect non-TTY → plain output (`notes list`), or exit 2 with a message.

**E7:** linux/amd64+arm64, darwin/amd64+arm64, windows/amd64(+arm64); GoReleaser on tag; checksums + signatures; smoke-test binaries in CI on each OS; distribute via GitHub Releases, Homebrew tap, Scoop, `go install`, container image.

**E8:**
```go
func writeAtomic(path string, data []byte) error {
    dir := filepath.Dir(path)
    f, err := os.CreateTemp(dir, ".tmp-*"); if err != nil { return err }
    defer os.Remove(f.Name())
    if _, err := f.Write(data); err != nil { f.Close(); return err }
    if err := f.Sync(); err != nil { f.Close(); return err }
    if err := f.Close(); err != nil { return err }
    return os.Rename(f.Name(), path)
}
```

**E9:** `View` runs only inside the running program, and startup failures happen before that; also errors displayed inside the alt-screen vanish on exit. Return the error from setup, quit, then print to stderr with a non-zero exit.

**E10:** 1 Windows Terminal; 2 macOS Terminal/iTerm2; 3 Linux GNOME Terminal; 4 tmux; 5 VS Code terminal; 6 SSH session; 7 80×24 and 40×10 and 300×80; 8 `NO_COLOR`/16-color/truecolor; 9 light and dark themes; 10 pipe/CI (non-TTY) behavior + `--json`.
