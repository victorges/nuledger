package rule

import (
	"nuledger/model"
)

// CommitFunc is a function that can be returned by rule.Authorizer's to be
// called as a confirmation that the transaction was executed. It can be used
// to update the internal state of the authorizer so as to guarantee future
// authorizations are performed consistently, considering only the actually
// executed transactions.
type CommitFunc func()

// Authorizer is the interface for any rule to be enforced when authorizing a
// transaction on a given account. If the rule is broken, an error should be
// returned by the Authorize function, which can be of type violation.Error to
// return the corresponding error code to the user (instead of a fatal error).
// It can also return a CommitFunc which should be called by the transaction
// execution agent in case the transaction is actually executed in the end
// (i.e. if no other rules blocked that transaction).
type Authorizer interface {
	Authorize(account model.Account, transaction model.Transaction) (CommitFunc, error)
}

// AuthorizerFunc is an adapter to allow the use of ordinary functions as
// authorizers of corresponding rules. If f is a function with the appropriate
// signature, AuthorizerFunc(f) is an Authorizer that calls f.
type AuthorizerFunc func(account model.Account, transaction model.Transaction) (CommitFunc, error)

// Authorize calls f(account, transaction) and returns its output.
func (f AuthorizerFunc) Authorize(account model.Account, transaction model.Transaction) (CommitFunc, error) {
	return f(account, transaction)
}
