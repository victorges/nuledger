// Package rules defines actual implementations of some rule authorizers. Should
// be used by a transaction performing agent for authorizing transactions on an
// account.
package rules

import (
	"nuledger/authorizer/rule"
	"nuledger/model"
	"nuledger/model/violation"
)

// AccountCardActive is a rule.AuthorizerFunc to check if the account card is
// active and returns a card-not-active violation error otherwise.
func AccountCardActive(account model.Account, _ model.Transaction) (rule.CommitFunc, error) {
	if !account.ActiveCard {
		return nil, violation.ErrorCardNotActive
	}
	return nil, nil
}

// SufficientLimit is a rule.AuthorizerFunc to check if the account has
// sufficient limit for performing the given transaction and returns an
// insufficient-limit violation error otherwise.
func SufficientLimit(account model.Account, transaction model.Transaction) (rule.CommitFunc, error) {
	if account.AvailableLimit < transaction.Amount {
		return nil, violation.ErrorInsufficientLimit
	}
	return nil, nil
}
