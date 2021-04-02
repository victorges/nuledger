package rules

import (
	"nuledger/authorizer/rule"
	"time"
)

const (
	frequencyAnalysisInterval = 2 * time.Minute
	maxIntervalTransactions   = 3
)

func Default() rule.Authorizer {
	return rule.List{
		&ChronologicalOrder{},
		rule.AuthorizerFunc(AccountCardActive),
		rule.AuthorizerFunc(SufficientLimit),
		NewLimitedFrequency(maxIntervalTransactions, frequencyAnalysisInterval),
		NewNoDoubleTransaction(frequencyAnalysisInterval),
	}
}
