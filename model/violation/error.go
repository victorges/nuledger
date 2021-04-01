package violation

import "fmt"

// Error is an error object to represent a violation error. It contains a
// validation code with a predictable message for the user.
type Error struct {
	Code
	Message string
}

// NewError creates a new violation error with the provided code and message.
func NewError(code Code, format string, args ...interface{}) error {
	return &Error{code, fmt.Sprintf(format, args...)}
}

func (e *Error) Error() string {
	return e.Message
}
