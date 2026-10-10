package main

import (
	"bytes"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/exp/teatest/v2"
)

func press(s string) tea.KeyPressMsg {
	switch s {
	case "up":
		return tea.KeyPressMsg{Code: tea.KeyUp}
	case "down":
		return tea.KeyPressMsg{Code: tea.KeyDown}
	}
	return tea.KeyPressMsg{Code: rune(s[0]), Text: s}
}

// Level 1: Update is a pure function - no terminal needed.
func TestUpdate_DownNeverBelowZero(t *testing.T) {
	m := model{}
	next, cmd := m.Update(press("down"))
	if got := next.(model).count; got != 0 {
		t.Fatalf("count = %d, want 0", got)
	}
	if cmd != nil {
		t.Fatal("expected no command")
	}
}

func TestUpdate_QuitReturnsQuitCmd(t *testing.T) {
	_, cmd := model{}.Update(press("q"))
	if cmd == nil {
		t.Fatal("expected a command")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatalf("cmd() = %T, want tea.QuitMsg", cmd())
	}
}

// Level 2: run the real program loop headlessly with teatest.
func TestProgram_EndToEnd(t *testing.T) {
	tm := teatest.NewTestModel(t, model{}, teatest.WithInitialTermSize(40, 10))
	tm.Send(press("up"))
	tm.Send(press("up"))
	teatest.WaitFor(t, tm.Output(), func(b []byte) bool {
		return bytes.Contains(b, []byte("count: 2"))
	}, teatest.WithDuration(2*time.Second))
	tm.Send(press("q"))
	final := tm.FinalModel(t, teatest.WithFinalTimeout(2*time.Second)).(model)
	if final.count != 2 || !final.done {
		t.Fatalf("final = %+v", final)
	}
}
