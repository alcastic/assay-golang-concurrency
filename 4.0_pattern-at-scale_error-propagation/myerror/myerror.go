package myerror

import "runtime/debug"

type MyError struct {
	Inner      error
	Message    string
	StackTrace string
}

func (err MyError) Error() string {
	return err.Message
}

func WrapError(err error, msg string) MyError {
	return MyError{
		Inner:      err,
		Message:    msg,
		StackTrace: string(debug.Stack()),
	}
}
