package rules

import "time"

func Default() []Rule {
	return []Rule{
		&ChronologicalOrder{},
		RuleFunc(AccountCardActive),
		RuleFunc(SufficientLimit),
		NewLimitedFrequency(3, 2*time.Minute),
		NewNoDoubleTransaction(2 * time.Minute),
	}
}
