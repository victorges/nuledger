package rule

import "nuledger/model"

// AuthorizerFunc is an adapter to use ordinary functions as rule authorizers.
// If f is a function with the appropriate signature, AuthorizerFunc(f) is an
// Authorizer that calls f.
type AuthorizerFunc func(account model.Account, transaction model.Transaction) (CommitFunc, error)

// Authorize calls f(account, transaction) and returns its output.
func (f AuthorizerFunc) Authorize(account model.Account, transaction model.Transaction) (CommitFunc, error) {
	return f(account, transaction)
}
