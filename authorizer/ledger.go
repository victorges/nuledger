package authorizer

import (
	"nuledger/authorizer/rule"
	"nuledger/model"
	"nuledger/model/violation"
)

type Ledger struct {
	accountState *model.Account
	rules        rule.List
}

func NewLedger(rules rule.List) *Ledger {
	return &Ledger{rules: rules}
}

func (l *Ledger) CreateAccount(account *model.Account) (*model.Account, error) {
	if l.accountState != nil {
		return l.accountState.Copy(), violation.ErrorAccountAlreadyInitialized
	}

	l.accountState = account.Copy()
	return account, nil
}

func (l *Ledger) PerformTransaction(transaction *model.Transaction) (*model.Account, error) {
	if l.accountState == nil {
		return nil, violation.ErrorAccountNotInitialized
	}
	account := l.accountState

	commitFunc, err := l.rules.Authorize(*account, transaction)
	if err != nil {
		return account.Copy(), err
	}

	account.AvailableLimit -= transaction.Amount
	commitFunc()
	return account.Copy(), nil
}
