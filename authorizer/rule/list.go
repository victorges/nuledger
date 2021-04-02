package rule

import "nuledger/model"

type List []Rule

// Ensure List implements of Rule interface
var _ Rule = List(nil)

func (l List) Authorize(account model.Account, transaction *model.Transaction) (CommitFunc, error) {
	var (
		commitFuncs = make([]CommitFunc, 0, 2)
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
	return combine(commitFuncs), model.AggregateErrors(errs)
}

func combine(funcs []CommitFunc) CommitFunc {
	return func() {
		for _, f := range funcs {
			f()
		}
	}
}
