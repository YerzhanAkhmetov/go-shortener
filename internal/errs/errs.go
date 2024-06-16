package errs

import (
	"fmt"
)

type Error struct {
	Msg    string `json:"message"`
	Code   int    `json:"code"`
	Status string `json:"status"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s\n%d\n%s", e.Msg, e.Code, e.Status)
}

func NewError(msg string, code int, status string) *Error {
	return &Error{
		Msg:    msg,
		Code:   code,
		Status: status,
	}
}
