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

//go:generate ../gen_mocks.sh ledger.go

// Ledger is the main component responsible for managing the account and
// performing transactions on it.
type Ledger interface {
	// CreateAccount creates a new account in the ledger.  It returns the final
	// state of the account managed by the ledger and any error that prevent it
	// from being created.
	//
	// The returned account has either the state of the newly created account,
	// the state of the existing account if one was already there, or nil if
	// some other type of error ocurred.
	CreateAccount(account model.Account) (*model.Account, error)
	// PerformTransaction receives a transaction and tries to perform it on the
	// managed account. It returns the final state of the account and any error
	// encountered that caused the operation to fail.
	//
	// PerformTransaction must return a non-nil error if the transaction was not
	// performed, in which case the returned account state must be the same
	// unmodified state of the account as before the attempt, with nil
	// representing a non-exsiting account. If the transaction is performed
	// successfully, the returned account will have the updated state (balance).
	PerformTransaction(transaction model.Transaction) (*model.Account, error)
}

// NewLedger creates an AuthLedger object with the provided Authorizer, which is
// used to authorize any attempt to perform a transaction.
//
// The only validations it performs itself are the ones regarding the actual
// creation and existence of the account. It is the sole one responsible for
// that since the Authorizers can only receive existing account and transaction
// objects by value (otherwise they'd all have to repeat the same not-nil
// validation themselves).
func NewLedger(authorizer rule.Authorizer) *AuthLedger {
	return &AuthLedger{authzer: authorizer}
}

// AuthLedger is the implementation of the Ledger interface delegating to a
// rule.Authorizer to authorize all the transactions. Not to be confused with
// Heath Ledger actor.
type AuthLedger struct {
	accountState *model.Account
	authzer      rule.Authorizer
}

// CreateAccount implements the Ledger interface. It currently only supports a
// single account, so this can be called only once per ledger instance or an
// account-already-initialized error will be returned.
func (l *AuthLedger) CreateAccount(account model.Account) (*model.Account, error) {
	if l.accountState != nil {
		return l.accountState.Copy(), violation.ErrorAccountAlreadyInitialized
	}

	l.accountState = &account
	return l.accountState.Copy(), nil
}

// PerformTransaction implements the Ledger interface. It initially calls the
// configured authorizer to ensure that the transaction is allowed and then
// performs it updating the current state of the account.
func (l *AuthLedger) PerformTransaction(transaction model.Transaction) (*model.Account, error) {
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
