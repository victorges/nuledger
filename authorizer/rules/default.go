package rules

import (
	"nuledger/authorizer/rule"
	"time"
)

const (
	frequencyAnalysisInterval = 2 * time.Minute
	maxIntervalTransactions   = 3
)

func Default() rule.RuleList {
	return rule.RuleList{
		&ChronologicalOrder{},
		rule.RuleFunc(AccountCardActive),
		rule.RuleFunc(SufficientLimit),
		NewLimitedFrequency(maxIntervalTransactions, frequencyAnalysisInterval),
		NewNoDoubleTransaction(frequencyAnalysisInterval),
	}
}
