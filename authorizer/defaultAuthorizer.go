package authorizer

import (
	"nuledger/authorizer/rule"
	"nuledger/authorizer/rules"
	"time"
)

const (
	frequencyAnalysisInterval = 2 * time.Minute
	maxIntervalTransactions   = 3
)

// DefaultAuthorizer returns an Authorizer with all the default rules to be
// validated for every transaction in the system. We could say that this gathers
// most of the core business logic validations that we want to perform against
// the transactions in order to authorize them or return their violations.
func DefaultAuthorizer() rule.Authorizer {
	return rule.List{
		&rules.ChronologicalOrder{},
		rule.AuthorizerFunc(rules.AccountCardActive),
		rule.AuthorizerFunc(rules.SufficientLimit),
		rules.NewLimitedFrequency(maxIntervalTransactions, frequencyAnalysisInterval),
		rules.NewUniqueTransactions(frequencyAnalysisInterval),
	}
}
