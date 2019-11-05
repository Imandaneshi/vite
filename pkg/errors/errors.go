package errors

import "fmt"


// Error is the main error objects used in both code and API error responses
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	source   error

}

func (e *Error) Error() string {
	return e.Code
}

func (e *Error) Source() string {
	if e.source != nil {
		return e.source.Error()
	}
	return ""
}

// generates a new error object
func New(c string, msg string, err error) error {
	return &Error{
		Code:    c,
		Message: msg,
		source:   err,
	}
}

// common errors

const alreadyExistsCode = "alreadyExists"
func AlreadyExistsError(msg string, err error, arg ...interface{}) error {
	return New(alreadyExistsCode, fmt.Sprintf(msg, arg...), err)
}

const validationErrorCode = "validationError"
func ValidationError(msg string, err error, arg ...interface{}) error {
	return New(validationErrorCode, fmt.Sprintf(msg, arg...), err)
}

const InternalServerErrorCode = "internalServerError"
func InternalServerError() error {
	return New(InternalServerErrorCode, "Internal server error", nil)
}

const NotFoundErrorCode = "notFound"
func NotFoundError(msg string) error {
	return New(NotFoundErrorCode, msg, nil)
}

const authenticationRequiredErrorCode = "authenticationRequired"
const authenticationRequiredErrorMessage = "Authentication is required"
func AuthenticationRequired() error {
	return New(authenticationRequiredErrorCode, authenticationRequiredErrorMessage, nil)
}