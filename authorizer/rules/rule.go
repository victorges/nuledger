package rules

import (
	"nuledger/model"
)

type CommitFunc func()

type Rule interface {
	Authorize(account model.Account, transaction *model.Transaction) (CommitFunc, error)
}

type RuleFunc func(account model.Account, transaction *model.Transaction) (CommitFunc, error)

func (f RuleFunc) Authorize(account model.Account, transaction *model.Transaction) (CommitFunc, error) {
	return f(account, transaction)
}
