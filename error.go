package kid

import "fmt"

type Code int

type Error struct {
	Status   *Status
	Original error
}

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Status.Message, e.Original.Error())
}

func NewError(status *Status, originals ...error) *Error {
	err := Error{Status: status}
	if len(originals) > 1 {
		err.Original = originals[0]
	}
	return &err
}
