package fsm

import "errors"

type (
	State      uint8
	Event      string
	Action     func(form, to State, args ...interface{}) error
	Transition struct {
		From   State
		To     State
		Event  Event
		Action Action
	}
)

type FSM struct {
	transitions []Transition
}

func NewFSM(transitions ...Transition) *FSM {
	return &FSM{
		transitions: transitions,
	}
}

func (f *FSM) Trigger(currentState State, event Event, args ...interface{}) error {
	trans := f.findTrans(currentState, event)
	if trans == nil {
		return errors.New("未找到Transition")
	}
	if err := trans.Action(trans.From, trans.To, args...); err != nil {
		return err
	}
	return nil
}

func (f *FSM) findTrans(form State, event Event) *Transition {
	for _, t := range f.transitions {
		if t.From == form && t.Event == event {
			return &t
		}
	}
	return nil
}
