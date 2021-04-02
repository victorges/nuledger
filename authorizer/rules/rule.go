package rules

import (
	"nuledger/model"
)

type CommitFunc func()

type Rule interface {
	Validate(account model.Account, transaction *model.Transaction) (CommitFunc, error)
}
