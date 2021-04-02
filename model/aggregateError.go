package model

import (
	"fmt"
	"strings"
)

type AggregateError struct {
	Errors []error
}

func AggregateErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
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
