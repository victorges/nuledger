package rule

import (
	"nuledger/model"
)

type CommitFunc func()

type Authorizer interface {
	Authorize(account model.Account, transaction *model.Transaction) (CommitFunc, error)
}

type AuthFunc func(account model.Account, transaction *model.Transaction) (CommitFunc, error)

func (f AuthFunc) Authorize(account model.Account, transaction *model.Transaction) (CommitFunc, error) {
	return f(account, transaction)
}
