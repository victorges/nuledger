package authorizer

import (
	"nuledger/authorizer/rule"
	"nuledger/model"
	"nuledger/model/violation"
)

type Authorizer struct {
	accountState *model.Account
	rules        rule.List
}

func NewAuthorizer(rules rule.List) *Authorizer {
	return &Authorizer{rules: rules}
}

func (a *Authorizer) CreateAccount(account *model.Account) (*model.Account, error) {
	if a.accountState != nil {
		return a.accountState.Copy(), violation.ErrorAccountAlreadyInitialized
	}

	a.accountState = account.Copy()
	return account, nil
}

func (a *Authorizer) PerformTransaction(transaction *model.Transaction) (*model.Account, error) {
	if a.accountState == nil {
		return nil, violation.ErrorAccountNotInitialized
	}
	account := a.accountState

	commitFunc, err := a.rules.Authorize(*account, transaction)
	if err != nil {
		return account.Copy(), err
	}

	account.AvailableLimit -= transaction.Amount
	commitFunc()
	return account.Copy(), nil
}
