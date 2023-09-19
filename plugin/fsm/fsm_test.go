package fsm

import (
	"fmt"
	"testing"
)

type Turnstile struct {
	ID         uint64
	EventCount uint64
	CoinCount  uint64
	PassCount  uint64
	State      State
}

const (
	StateOpen = 1
	StateLock = 2
)

func Coin(form, to State, args ...interface{}) error {
	ts := args[0].(*Turnstile)
	fmt.Printf("转门 [%d] 投币\n", ts.ID)
	ts.CoinCount++
	ts.State = to
	return nil
}

// Push 推门
func Push(form, to State, args ...interface{}) error {
	ts := args[0].(*Turnstile)
	fmt.Printf("转门 [%d] 打开, 通过\n", ts.ID)
	ts.PassCount++
	ts.State = to
	return nil
}

func Opened(form, to State, args ...interface{}) error {
	ts := args[0].(*Turnstile)
	fmt.Printf("转门 [%d] 已打开, 请勿重复投币\n", ts.ID)
	return nil
}

func Locked(form, to State, args ...interface{}) error {
	ts := args[0].(*Turnstile)
	fmt.Printf("转门 [%d] 已锁, 无法通过\n", ts.ID)
	return nil
}

func TestFSM(t *testing.T) {
	ts := &Turnstile{
		ID:    1,
		State: StateLock,
	}
	fsm := initFSM()

	// 直接推门
	if err := fsm.Trigger(ts.State, "Pass", ts); err != nil {
		t.Errorf("trigger err: %v", err)
	}

	// 多次投币
	if err := fsm.Trigger(ts.State, "Coin", ts); err != nil {
		t.Errorf("trigger err: %v", err)
	}
	if err := fsm.Trigger(ts.State, "Coin", ts); err != nil {
		t.Errorf("trigger err: %v", err)
	}

	// 先投币再多次推门
	if err := fsm.Trigger(ts.State, "Coin", ts); err != nil {
		t.Errorf("trigger err: %v", err)
	}
	if err := fsm.Trigger(ts.State, "Pass", ts); err != nil {
		t.Errorf("trigger err: %v", err)
	}
	if err := fsm.Trigger(ts.State, "Pass", ts); err != nil {
		t.Errorf("trigger err: %v", err)
	}

	// 正常流程 先投币再推门
	if err := fsm.Trigger(ts.State, "Coin", ts); err != nil {
		t.Errorf("trigger err: %v", err)
	}
	if err := fsm.Trigger(ts.State, "Pass", ts); err != nil {
		t.Errorf("trigger err: %v", err)
	}
	if err := fsm.Trigger(ts.State, "Coin", ts); err != nil {
		t.Errorf("trigger err: %v", err)
	}
	if err := fsm.Trigger(ts.State, "Pass", ts); err != nil {
		t.Errorf("trigger err: %v", err)
	}
}

func initFSM() *FSM {
	fsm := NewFSM(
		Transition{From: StateLock, To: StateOpen, Event: "Coin", Action: Coin},
		Transition{From: StateOpen, To: StateLock, Event: "Pass", Action: Push},
		Transition{From: StateOpen, To: StateOpen, Event: "Coin", Action: Opened},
		Transition{From: StateLock, To: StateLock, Event: "Pass", Action: Locked},
	)
	return fsm
}
