package rule

import (
	"nuledger/model"
)

type CommitFunc func()

type Rule interface {
	Authorize(account model.Account, transaction *model.Transaction) (CommitFunc, error)
}

type Func func(account model.Account, transaction *model.Transaction) (CommitFunc, error)

func (f Func) Authorize(account model.Account, transaction *model.Transaction) (CommitFunc, error) {
	return f(account, transaction)
}
