# Unit 14 Quiz

**Q1 (easy).** Which is the conventional exit code for a usage/config error?
A) 0  B) 1  C) 2  D) 255

**Q2 (easy).** In Wish, how many models per SSH connection?
A) One shared for all  B) One per session  C) One per keypress  D) One per server

**Q3 (medium).** Where should errors be printed after a TUI fails?
A) stdout during the TUI  B) stderr after `Run` returns  C) Inside `View`  D) They shouldn't

**Q4 (medium).** Best place for domain logic:
A) `Update`  B) A package without Bubble Tea imports  C) `main`  D) `View`

**Q5 (medium).** How do you cross-compile for Linux ARM64?
A) Install a cross toolchain  B) `GOOS=linux GOARCH=arm64 go build`  C) `go cross`  D) Docker only

**Q6 (medium).** Short answer: why check `term.IsTerminal` before starting a TUI?

**Q7 (hard).** Short answer: a live-update goroutine leaks per SSH session. Cause and fix.

**Q8 (medium).** Which is the best config precedence?
A) file > flags > env  B) defaults < file < env < flags  C) flags < env < file  D) env only

**Q9 (hard, apply).** Write the `-ldflags` and variable to embed a version string.

**Q10 (medium).** Give two things Wish does *not* do for your security.
