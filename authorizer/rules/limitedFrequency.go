package rules

import (
	"nuledger/authorizer/rule"
	"nuledger/authorizer/util"
	"nuledger/model"
	"nuledger/model/violation"
	"time"
)

type LimitedFrequency struct {
	limiter *util.RateLimiter
}

func NewLimitedFrequency(maxTransactions int, interval time.Duration) rule.Authorizer {
	return &LimitedFrequency{util.NewRateLimiter(maxTransactions, interval)}
}

func (f *LimitedFrequency) Authorize(_ model.Account, transaction *model.Transaction) (rule.CommitFunc, error) {
	if !f.limiter.Allows(transaction.Time) {
		return nil, violation.ErrorHighFrequencySmallInterval
	}
	commit := func() { f.limiter.Take(transaction.Time) }
	return commit, nil
}
