package rules

import "time"

const (
	frequencyAnalysisInterval = 2 * time.Minute
	maxIntervalTransactions   = 3
)

func Default() []Rule {
	return []Rule{
		&ChronologicalOrder{},
		RuleFunc(AccountCardActive),
		RuleFunc(SufficientLimit),
		NewLimitedFrequency(maxIntervalTransactions, frequencyAnalysisInterval),
		NewNoDoubleTransaction(frequencyAnalysisInterval),
	}
}
