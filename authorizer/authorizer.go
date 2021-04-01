package authorizer

import (
	"nuledger/model"
)

type Authorizer struct{}

func (a *Authorizer) CreateAccount(account *model.Account) (model.Account, error) {
	return *account, nil
}

func (a *Authorizer) PerformTransaction(transaction *model.Transaction) (model.Account, error) {
	return model.Account{}, nil
}
