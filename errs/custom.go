package errs

import "fmt"

type NothingToDo struct{}

func (NothingToDo) Error() string {
	return "nothing to do"
}

func IsNothingToDo(err error) bool {
	_, ok := err.(NothingToDo)
	return ok
}

type NothingFound struct{}

func (NothingFound) Error() string {
	return "not found"
}

func IsNotFound(err error) bool {
	_, ok := err.(NothingFound)
	return ok
}

type StateErr struct {
	s   uint8
	msg string
}

func WithState(state uint8, msg string) *StateErr {
	return &StateErr{s: state, msg: msg}
}

func (e StateErr) Error() string {
	return fmt.Sprintf("msg: %s [[state: %d]]", e.msg, e.s)
}

func InState(err error, state uint8) bool {
	e, ok := err.(*StateErr)
	if !ok {
		return false
	}
	return e.s == state
}
