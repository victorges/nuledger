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

// Ledger is the main component responsible for managing the account and
// performing transactions on it. It receives an Authorizer on creation and uses
// it to authorize any transaction that may be performed.
//
// The only validations it performs itself are the ones regarding the actual
// creation and existence of the account. It is the sole one responsible for
// that and the Authorizers can only receive existing account and transaction
// objects by value (otherwise they'd all have to repeat the same not-nil
// validation themselves).
type Ledger struct {
	accountState *model.Account
	authzer      rule.Authorizer
}

// NewLedger creates a ledger object with the provided authorizer.
func NewLedger(authorizer rule.Authorizer) *Ledger {
	return &Ledger{authzer: authorizer}
}

// CreateAccount creates a new account in the ledger. It currently only supports
// a single account, so this can be called only once per ledger instance or an
// account-already-initialized error will be returned. It also returns the final
// state of the account managed by the ledger, either the state of the newly
// created account or the state of the existing one in case one had already been
// created.
func (l *Ledger) CreateAccount(account model.Account) (*model.Account, error) {
	if l.accountState != nil {
		return l.accountState.Copy(), violation.ErrorAccountAlreadyInitialized
	}

	l.accountState = &account
	return l.accountState.Copy(), nil
}

// PerformTransaction receives a transaction and tries to perform it on the
// managed account. It initially calls the configured authorizer to ensure that
// the transaction is allowed and then performs it updating the current state of
// the account. It returns the final state of the account, either with the
// updated balance if the transaction was performed or the same state as before
// in case it was not authorized.
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
