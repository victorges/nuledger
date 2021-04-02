package rules

import (
	"nuledger/model"
	"nuledger/model/violation"
	"time"
)

type LimitedFrequency struct {
	limiter *RateLimiter
}

func NewLimitedFrequency(maxTransactions int, interval time.Duration) Rule {
	return &LimitedFrequency{NewRateLimiter(maxTransactions, interval)}
}

func (f *LimitedFrequency) Authorize(_ model.Account, transaction *model.Transaction) (CommitFunc, error) {
	if !f.limiter.Allows(transaction.Time) {
		return nil, violation.ErrorHighFrequencySmallInterval
	}
	commit := func() { f.limiter.Take(transaction.Time) }
	return commit, nil
}
