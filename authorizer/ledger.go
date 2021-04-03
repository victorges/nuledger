// Package authorizer contains the higher-level components of the authorizer
// logic.
//
// It's main components are a Ledger which handles the core business
// logic of creating accounts and performing transactions and a Handler which
// translates between lower-level JSON messages and the Ledger interface.
package authorizer

import (
	"nuledger/authorizer/rule"
	"nuledger/model"
	"nuledger/model/violation"
)

type Ledger struct {
	accountState *model.Account
	authzer      rule.Authorizer
}

func NewLedger(authzer rule.Authorizer) *Ledger {
	return &Ledger{authzer: authzer}
}

func (l *Ledger) CreateAccount(account model.Account) (*model.Account, error) {
	if l.accountState != nil {
		return l.accountState.Copy(), violation.ErrorAccountAlreadyInitialized
	}

	l.accountState = &account
	return l.accountState.Copy(), nil
}

func (l *Ledger) PerformTransaction(transaction model.Transaction) (*model.Account, error) {
	if l.accountState == nil {
		return nil, violation.ErrorAccountNotInitialized
	}
	account := l.accountState

	commitFunc, err := l.authzer.Authorize(*account, transaction)
	if err != nil {
		return account.Copy(), err
	}

	account.AvailableLimit -= transaction.Amount
	commitFunc()
	return account.Copy(), nil
}
