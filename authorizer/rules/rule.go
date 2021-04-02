package rules

import "nuledger/model"

type CommitFunc func()

type Rule interface {
	Validate(transaction *model.Transaction) (CommitFunc, error)
}

func Default() []Rule {
	return []Rule{
		&ChronologicalOrder{},
	}
}
