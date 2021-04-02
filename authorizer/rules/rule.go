package rules

import (
	"nuledger/model"
)

type CommitFunc func()

type Rule interface {
	Authorize(account model.Account, transaction *model.Transaction) (CommitFunc, error)
}
