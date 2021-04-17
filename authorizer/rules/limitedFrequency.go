package rules

import (
	"nuledger/authorizer/rule"
	"nuledger/model"
	"nuledger/model/violation"
	"nuledger/util"
	"time"
)

// NewLimitedFrequency returns a rule.Authorizer to guarantee that transactions
// do not exceed a maximum allowed frequency. That frequency can be configured
// through the arguments passed to this constructor, so that at most
// `maxTransactions` are authorized under the given `interval`.
//
// Its Authorize function checks if the maximum allowed frequency is exceeded,
// and if so the transaction is not authorized and a violation error of
// high-frequency-small-interval is returned.
func NewLimitedFrequency(maxTransactions int, interval time.Duration) rule.Authorizer {
	limiter := util.RateLimiter{MaxEvents: maxTransactions, Interval: interval}
	keyMapper := func(tx *model.Transaction) interface{} {
		return tx.AccountID
	}
	return NewFrequencyAnalyzer(limiter, keyMapper, violation.ErrorHighFrequencySmallInterval)
}
