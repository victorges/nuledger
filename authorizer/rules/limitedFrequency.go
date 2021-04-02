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

func (f *LimitedFrequency) Validate(_ model.Account, transaction *model.Transaction) (CommitFunc, error) {
	if !f.limiter.Allow(transaction.Time) {
		err := violation.NewError(violation.HighFrequencySmallInterval, "Too many transactions in a small interval")
		return nil, err
	}
	commit := func() { f.limiter.Take(transaction.Time) }
	return commit, nil
}
