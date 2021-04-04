package util

import (
	"fmt"
	"strings"
)

// AggregateError is a type for representing multiple errors as a single one.
// Meant for multiplexed or even parallelized operations that may have multiple
// errors to return while consumers still may need to process each of the errors
// separately.
type AggregateError struct {
	Errors []error
}

// AggregateErrors receives a slice of errors and returns an appropriate
// representation of them as a single error. For an empty slice, it returns a
// nil error; for a slice with only 1 element, it returns that single error; and
// finally for bigger slices, it returns an AggregateError object containing all
// the errors in the slice.
func AggregateErrors(errs []error) error {
	if count := len(errs); count == 0 {
		return nil
	} else if count == 1 {
		return errs[0]
	}
	return AggregateError{errs}
}

// Error implements the error interface. It combines the error messages of all
// the aggregated errors and returns a single string representing them.
func (a AggregateError) Error() string {
	msgs := make([]string, len(a.Errors))
	for i, err := range a.Errors {
		msgs[i] = err.Error()
	}
	return fmt.Sprintf("Multiple errors: %s", strings.Join(msgs, "; "))
}
