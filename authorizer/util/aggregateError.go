package util

import (
	"fmt"
	"strings"
)

type AggregateError struct {
	Errors []error
}

func AggregateErrors(errs []error) error {
	if count := len(errs); count == 0 {
		return nil
	} else if count == 1 {
		return errs[0]
	}
	return AggregateError{errs}
}

func (a AggregateError) Error() string {
	msgs := make([]string, len(a.Errors))
	for i, err := range a.Errors {
		msgs[i] = err.Error()
	}
	return fmt.Sprintf("Multiple errors: %s", strings.Join(msgs, "; "))
}
