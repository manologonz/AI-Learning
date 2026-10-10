# Unit 6 Exam

**E1 (Apply).** Write styles and a `View` for a two-column layout: left 30 columns fixed (bordered, no margin), right takes the rest (bordered), both full height minus a 1-line footer. (Remember what `Width` includes in v2.)

**E2 (Analyze).** Odd terminal width 81: `left := s.Width(m.w/2)`, `right := s.Width(m.w/2)`. What's wrong, and the fix?

**E3 (Analyze).** With numbers: for `Border(NormalBorder()).Padding(1,2).Width(20)`, what is (a) the rendered block width, (b) the content width available for text? How does adding `.Margin(0,3)` change (a)?

**E4 (Apply).** Choose a text color for both light and dark backgrounds using the v2 API, including the messages you must handle.

**E5 (Evaluate).** A colleague computes `strings.Repeat(" ", w-len(label))` to right-align styled text. Critique and give the corrected approach.

**E6 (Analyze).** Why should styles be defined once at package level (or in the model constructor) rather than inside `View`?

**E7 (Apply).** Build a status bar: left text, right text, both on one line filling `m.w`, with the gap filled by spaces.

**E8 (Evaluate).** Static output (e.g. `mytool --list`) needs a colored table. Which Lip Gloss subpackage vs Bubbles component, and why?

**E9 (Analyze).** Why is `lipgloss.Println` needed outside Bubble Tea but not inside?

**E10 (Create).** Design a responsive layout policy: ≥100 cols → sidebar + main; <100 → main only with the sidebar toggled by a key. Describe model fields and `View` branching, including tiny sizes (<20 cols).
