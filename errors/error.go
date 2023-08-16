package errors

import (
	"errors"
	"fmt"
)

type Code int

type Error struct {
	Status   *Status
	Original error
}

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var UnknownError = &Status{Code: 10000, Message: "未知错误"}

func (status *Status) FmtMsg(val ...interface{}) *Status {
	return &Status{
		Code:    status.Code,
		Message: fmt.Sprintf(status.Message, val...),
	}
}

func (e *Error) Error() string {
	var errStr string
	if e.Status != nil {
		errStr = fmt.Sprintf("code: %d message: %s", e.Status.Code, e.Status.Message)
	}
	if e.Original != nil {
		errStr = fmt.Sprintf("%s original error: %s", errStr, e.Original.Error())
	}
	return errStr
}

func New(status *Status, originals ...error) *Error {
	err := Error{Status: status}
	if len(originals) > 0 {
		err.Original = originals[0]
	}
	return &err
}

func (e *Error) GetStatus() *Status {
	return e.Status
}

func GetStatus(e error) *Status {
	var err *Error
	switch {
	case errors.As(e, &err):
		return err.GetStatus()
	default:
		return UnknownError
	}
}

func GetCodeMessage(e error) (int, string) {
	status := GetStatus(e)
	return status.Code, status.Message
}
