package errors

import "fmt"

type Error struct {
	Code    string
	Message string
	source   error
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Source() string {
	if e.source != nil {
		return e.source.Error()
	}
	return ""
}

func New(c string, msg string, err error) error {
	return &Error{
		Code:    c,
		Message: msg,
		source:   err,
	}
}

// common errors

const alreadyExistsCode = "already_exists"
func AlreadyExistsError(msg string, err error, arg ...interface{}) error {
	return New(alreadyExistsCode, fmt.Sprintf(msg, arg...), err)
}