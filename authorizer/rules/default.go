package rules

import (
	"nuledger/authorizer/rule"
	"time"
)

const (
	frequencyAnalysisInterval = 2 * time.Minute
	maxIntervalTransactions   = 3
)

func Default() rule.List {
	return rule.List{
		&ChronologicalOrder{},
		rule.AuthFunc(AccountCardActive),
		rule.AuthFunc(SufficientLimit),
		NewLimitedFrequency(maxIntervalTransactions, frequencyAnalysisInterval),
		NewNoDoubleTransaction(frequencyAnalysisInterval),
	}
}
