# Unit 15 Quiz — Answers

**Q1: B.** `core` stays UI-free for tests and reuse.

**Q2: B.** Without the ID, an old timer would clear a newer message early.

**Q3: B.** Route keys to the child while it's capturing input; check `FilterState() == list.Filtering` before global shortcuts.

**Q4: B.** Atomic replace prevents corrupted files on crash.

**Q5:** `loading`, `error` (with retry), `empty`, `ready`.

**Q6: B.** `Sequence` runs in order. `Batch` runs concurrently and could quit first.

**Q7:** Filtering/sorting/refresh changes indexes; the ID identifies the task stably.

**Q8:** In v2 `Width` includes border and padding (margin is extra); content width = `Width − GetHorizontalFrameSize()` (Unit 6).
