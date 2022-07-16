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

func (e *Error) GetStatus() (int, string) {
	return e.Status.Code, e.Status.Message
}

func Value(e error) (int, string) {
	switch err := e.(type) {
	case *Error:
		return err.GetStatus()
	default:
		return 500, "未知错误"
	}
}
