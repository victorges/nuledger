package rules

import (
	"nuledger/authorizer/rule"
	"nuledger/authorizer/util"
	"nuledger/model"
	"nuledger/model/violation"
	"time"
)

// LimitedFrequency is a rule.Authorizer to guarantee that transactions do not
// exceed a maximum allowed frequency, which can be configured on its creation.
type LimitedFrequency struct {
	limiter util.RateLimiter
}

// NewLimitedFrequency returns a LimitedFrequency authorizer with the given
// configuration. The arguments of the constructor configures it so that at most
// `maxTransactions` are performed under the given `interval`.
func NewLimitedFrequency(maxTransactions int, interval time.Duration) *LimitedFrequency {
	return &LimitedFrequency{
		limiter: util.RateLimiter{
			MaxEvents: maxTransactions,
			Interval:  interval,
		}}
}

// Authorize checks if the maximum allowed frequency is exceeded, and if so the
// transaction is not authorized and a high-frequency-small-interval violation
// error is returned.
func (f *LimitedFrequency) Authorize(_ model.Account, transaction model.Transaction) (rule.CommitFunc, error) {
	if !f.limiter.Allows(transaction.Time) {
		return nil, violation.ErrorHighFrequencySmallInterval
	}
	commit := func() { f.limiter.Take(transaction.Time) }
	return commit, nil
}
