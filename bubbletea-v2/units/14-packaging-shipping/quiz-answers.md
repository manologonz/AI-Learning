# Unit 14 Quiz — Answers

**Q1: C.** 2 by convention (0 success, 1 runtime failure).

**Q2: B.** Sessions are isolated; sharing a model would mix users' state.

**Q3: B.** After `Run` returns, the terminal is restored so stderr is visible and clean.

**Q4: B.** Keeps it testable and reusable in non-TUI/SSH modes.

**Q5: B.** Go cross-compiles natively (no cgo).

**Q6:** Without a TTY (pipes/CI) raw mode fails or escape codes pollute output; detect it and fall back to plain/JSON output or a clear error.

**Q7:** The goroutine uses `context.Background()`, so it outlives the disconnected session. Use `sess.Context()` (cancelled on disconnect) and select on `ctx.Done()`.

**Q8: B.** Later sources override earlier: defaults, then file, then env, then flags.

**Q9:** `var version = "dev"` in `main`; `go build -ldflags "-X main.version=v1.2.3"` (package path is `main`).

**Q10:** Input validation (paths, commands, sizes); authorization for what an authenticated user may do; rate limiting/resource caps; not running as root; secrets management.
