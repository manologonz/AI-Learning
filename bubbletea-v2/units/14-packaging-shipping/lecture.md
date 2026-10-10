# Unit 14: Packaging & Shipping

> Runnable code: `examples/u14-ssh/main.go` — a Wish SSH server hosting a Bubble Tea app, plus a `local` mode with TTY detection. (Run: `go run ./examples/u14-ssh local` or `go run ./examples/u14-ssh` then `ssh -p 23234 localhost`.)

## Objectives
By the end of this unit you can:
- Structure a TUI project so the model is a package separate from `main`.
- Add a CLI (flags/subcommands), a config file, and sensible exit codes.
- Detect non-interactive environments and provide fallbacks (`--json`, plain output).
- Serve a Bubble Tea app over SSH with Wish (per-session models, lifecycle).
- Build, version, and release cross-platform binaries.

## Prerequisites
Units 2–13.

## Core content

### Project layout
```
mytool/
  cmd/mytool/main.go      # thin: flags, wiring, tea.NewProgram, exit codes
  internal/ui/            # Bubble Tea models (app router, screens, components)
  internal/core/          # domain logic: no Bubble Tea imports; unit-testable
  internal/storage/       # files, DB, network clients (behind interfaces)
  go.mod
```
Rule: **`core` doesn't import Bubble Tea.** The UI calls core; core never knows about messages. This keeps logic testable (Unit 12) and reusable for a non-TUI mode and an SSH mode.

### `main` responsibilities
```go
func main() {
    cfg, err := loadConfig()               // flags/env/file
    if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(2) }   // usage/config error → 2
    p := tea.NewProgram(ui.New(cfg))
    if _, err := p.Run(); err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
}
```
Convention: `0` success, `1` runtime failure, `2` usage/config error. Print errors to **stderr** *after* the program returned (the terminal is restored by then).

### Flags & subcommands
Use the standard `flag` package for simple tools, or `spf13/cobra`/`urfave/cli` for subcommands (`mytool`, `mytool list --json`, `mytool serve`). Keep the TUI as the default command and add scriptable subcommands: users love a TUI that also has `--json`.

### Interactive vs. non-interactive
A TUI needs a TTY on stdin and stdout. Detect it (this course's example uses `github.com/charmbracelet/x/term`):
```go
if !term.IsTerminal(os.Stdin.Fd()) || !term.IsTerminal(os.Stdout.Fd()) {
    // fallback: print plain text / JSON, or return a clear error with exit code 2
}
```
Also: respect `NO_COLOR`, provide `--no-color`, and never require the TUI for automation. In Bubble Tea, `tea.WithoutRenderer()` lets the same program run as a plain CLI when output isn't a TTY; often simpler to branch in `main`.

### Config, state, and paths
- Config: `os.UserConfigDir()` → `~/.config/mytool/config.toml` (Linux), `%AppData%` on Windows; state/cache: `os.UserCacheDir()`.
- Precedence: defaults < file < env < flags.
- Write files atomically (temp file + rename) so a crash doesn't corrupt state.
- Never store secrets in plain text; use the OS keychain or env.

### Serving over SSH with Wish
Wish (`charm.land/wish/v2`) makes an SSH server whose sessions each get a Bubble Tea program, so users just `ssh yourhost`:
```go
s, err := wish.NewServer(
    wish.WithAddress(net.JoinHostPort("0.0.0.0", "23234")),
    wish.WithHostKeyPath(".ssh/id_ed25519"),          // generated if missing
    wish.WithMiddleware(
        wishtea.Middleware(func(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
            return ui.New(cfgFor(sess.User())), wishtea.MakeOptions(sess)   // per-session model + PTY IO
        }),
    ),
)
```
Key ideas:
- **One model per session**: never share a model between users; shared state (a chat room, a DB) lives outside and communicates via `Program.Send` or channels (Unit 10).
- `MakeOptions(sess)` wires the session's PTY as the program's input/output (and window size).
- **Per-session lifetime**: tie background goroutines to the session's context (`sess.Context()`), or you'll leak per connection (Unit 10 edge case).
- Add middleware for auth (public keys), logging, rate limiting, `activeterm` (require a PTY). **Order matters:** `wish.WithMiddleware(a, b, c)` composes first-to-last, so `c` runs *first*; list the app handler (`wishtea.Middleware`) first and cross-cutting ones (`activeterm`, `logging`) after it.
- Graceful shutdown: `s.Shutdown(ctx)` on SIGINT/SIGTERM (see the example).
- **Security**: SSH gives you an authenticated channel; you still must not expose a shell or arbitrary file access. Treat every key press as untrusted input; validate paths/commands. Don't run as root; use a dedicated user/container.

### Building and releasing
```bash
go build -trimpath -ldflags "-s -w -X main.version=$(git describe --tags)" ./cmd/mytool
GOOS=linux GOARCH=arm64 go build ...      # cross-compile: pure Go, no cgo needed
```
- Embed the version with `-ldflags -X`; add `--version`.
- **GoReleaser** automates cross-platform archives, checksums, Homebrew/Scoop taps, and GitHub releases from a tag.
- `go install example.com/mytool/cmd/mytool@latest` works for Go users.
- Consider Windows: Bubble Tea v2 supports it; test key handling and colors in Windows Terminal, and don't assume `$EDITOR`/`vi`.

### Compatibility matrix (test what you can)
Terminals: Windows Terminal, iTerm2, Kitty/WezTerm/Ghostty, GNOME Terminal, tmux/screen, VS Code terminal, over SSH. Color: truecolor / 256 / 16 / none (`NO_COLOR`). Light and dark backgrounds. Sizes: 80×24 and tiny/huge. Bracketed paste/mouse/focus availability varies (Unit 9): degrade gracefully.

### Documentation & help
Ship a README with a GIF/screencast (VHS by Charm records terminal GIFs from a script), a keybinding table (generated from your `key.Binding`s!), and in-app help (`?`).

## Worked examples

### Straightforward: version flag
```go
var version = "dev"
func main() { if len(os.Args) > 1 && os.Args[1] == "--version" { fmt.Println(version); return } … }
```
Build: `go build -ldflags "-X main.version=v1.2.3"`.

### Typical: TUI by default, `--json` for scripts
`mytool` → TUI. `mytool list --json` → calls `core.List()`, marshals, prints to stdout, exits: no Bubble Tea imported on that path (fast start, works in cron).

### Edge case: a TUI in a pipe, and per-session goroutine leaks
`mytool | less` starts a TUI with stdout being a pipe: escape codes pollute the output or the program errors on raw-mode setup. Detect and fall back (or exit 2 with "needs a terminal"). Over SSH, a session that disconnects abruptly ends its program, but a goroutine started with `context.Background()` for "live updates" never stops: use `sess.Context()`, and verify with `runtime.NumGoroutine()` after connect/disconnect cycles in a test.

## Common misconceptions
- **"`main` should contain the app."** Keep it thin; logic lives in packages.
- **"A TUI should always start."** Detect non-TTY and behave.
- **"One shared model for all SSH users."** One per session.
- **"SSH makes it secure."** It authenticates; your app still validates input.
- **"Cross-compiling needs special toolchains."** Pure Go: set GOOS/GOARCH.
- **"Errors print anywhere."** Print to stderr after `Run` returns.

## Check your understanding
1. Where should domain logic live? *(A package with no Bubble Tea imports.)*
2. What does `wishtea.MakeOptions(sess)` do? *(Wires the session's PTY IO/size to the program.)*
3. What exit code for bad flags? *(2.)*
4. How do you embed the version? *(`-ldflags "-X main.version=…"`.)*

## Summary
- Thin `main`, UI package, dependency-free core.
- TTY detection + non-interactive modes; config precedence; exit codes.
- Wish: a model per SSH session, per-session lifetimes, graceful shutdown.
- Cross-compile with `GOOS/GOARCH`; automate releases; test across terminals.
