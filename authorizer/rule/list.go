package rule

import (
	"nuledger/model"
	"nuledger/util"
)

// List is a helper type to allow the use of a slice of Authorizer objects as if
// it were a single Authorizer.
type List []Authorizer

// Ensure List implements the Authorizer interface
var _ Authorizer = List(nil)

// Authorize function from List type calls every authorizer in the slice and
// combines both the returned commit functions into a single CommitFunc which
// calls all of the original ones, and the returned errors into a possible
// model.AggregateError in case multiple errors were returned.
func (l List) Authorize(account model.Account, transaction model.Transaction) (CommitFunc, error) {
	var (
		commitFuncs = make([]CommitFunc, 0, 5)
		errs        []error
	)
	for _, rule := range l {
		commit, err := rule.Authorize(account, transaction)
		if commit != nil {
			commitFuncs = append(commitFuncs, commit)
		}
		if err != nil {
			errs = append(errs, err)
		}
	}
	return combine(commitFuncs), util.AggregateErrors(errs)
}

func combine(funcs []CommitFunc) CommitFunc {
	return func() {
		for _, f := range funcs {
			f()
		}
	}
}
