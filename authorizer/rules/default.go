package rules

import (
	"nuledger/authorizer/rule"
	"time"
)

const (
	frequencyAnalysisInterval = 2 * time.Minute
	maxIntervalTransactions   = 3
)

// Default returns an Authorizer with all the default rules expected to be
// validated for every transaction in the system. We could say that this gathers
// all the core business logic validations that we want to perform against the
// transactions in order to authorize them or return their violations.
func Default() rule.Authorizer {
	return rule.List{
		&ChronologicalOrder{},
		rule.AuthorizerFunc(AccountCardActive),
		rule.AuthorizerFunc(SufficientLimit),
		NewLimitedFrequency(maxIntervalTransactions, frequencyAnalysisInterval),
		NewUniqueTransactions(frequencyAnalysisInterval),
	}
}
