package violation

import "fmt"

// Error is an error object to represent a violation error.
type Error struct {
	// A well-known validation code for the user.
	Code
	// A free-form message that can contain a friendly description of the
	// violation and/or any additional info from the context.
	Message string
}

// NewError creates a new violation error with the provided code and message.
func NewError(code Code, format string, args ...interface{}) Error {
	return Error{code, fmt.Sprintf(format, args...)}
}

// Error implements the error interface to return the error message as a string.
func (e Error) Error() string {
	return e.Message
}
