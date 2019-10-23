package errors

import "fmt"

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
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

const validationErrorCode = "validation_error"
func ValidationError(msg string, err error, arg ...interface{}) error {
	return New(validationErrorCode, fmt.Sprintf(msg, arg...), err)
}

const InternalServerErrorCode = "internal_server_error"
func InternalServerError() error {
	return New(InternalServerErrorCode, "Internal server error", nil)
}