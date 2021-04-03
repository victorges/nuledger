package rule

import (
	"nuledger/model"
)

// An Authorizer enforces a rule when performing a transaction on an account.
//
// Authorize should return an error if the rule is broken, which can be of type
// violation.Error so the corresponding error code is returned to the user
// instead of being handled as a fatal error.
//
// It can also return a CommitFunc which will be called by the transaction
// execution agent in case the transaction is actually executed in the end
// i.e. if no other rules blocked that transaction.
type Authorizer interface {
	Authorize(account model.Account, transaction model.Transaction) (CommitFunc, error)
}

// CommitFunc is a function that can be returned by an Authorizer, for it to be
// called as a confirmation that the transaction was executed. It can be used
// to update the internal state of the Authorizer so as to guarantee future
// authorizations are performed consistently, considering only the actually
// executed transactions, not the attempted ones.
type CommitFunc func()

// AuthorizerFunc is an adapter to use ordinary functions as rule authorizers.
// If f is a function with the appropriate signature, AuthorizerFunc(f) is an
// Authorizer that calls f.
type AuthorizerFunc func(account model.Account, transaction model.Transaction) (CommitFunc, error)

// Authorize calls f(account, transaction) and returns its output.
func (f AuthorizerFunc) Authorize(account model.Account, transaction model.Transaction) (CommitFunc, error) {
	return f(account, transaction)
}
