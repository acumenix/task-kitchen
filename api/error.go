package api

import "fmt"

type errorType int

type userError struct {
	code  int
	msg   string
	cause error
}

func newUserError(code int, msg string, args ...interface{}) *userError {
	return &userError{
		code: code,
		msg:  fmt.Sprintf(msg, args...),
	}
}

func (x *userError) Error() string {
	return x.msg
}

func (x *userError) setCause(err error) *userError {
	x.cause = err
	return x
}
