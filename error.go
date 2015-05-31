package kempclient

import "fmt"

type Error struct {
	s    string
	Code int
}

func newError(errorCode int, cause string) *Error {
	return &Error{
		s:    cause,
		Code: errorCode,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.s)
}
