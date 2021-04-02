package rules

import (
	"nuledger/authorizer/rule"
	"nuledger/authorizer/util"
	"nuledger/model"
	"nuledger/model/violation"
	"time"
)

type NoDoubleTransaction struct {
	doubleTxInterval time.Duration
	limiters         map[doubleTransactionKey]*util.RateLimiter
}

func NewNoDoubleTransaction(interval time.Duration) rule.Authorizer {
	return &NoDoubleTransaction{interval, map[doubleTransactionKey]*util.RateLimiter{}}
}

func (d *NoDoubleTransaction) Authorize(_ model.Account, transaction model.Transaction) (rule.CommitFunc, error) {
	limiter := d.getLimiter(&transaction)
	if !limiter.Allows(transaction.Time) {
		return nil, violation.ErrorDoubleTransaction
	}
	commit := func() { limiter.Take(transaction.Time) }
	return commit, nil
}

type doubleTransactionKey struct {
	Merchant string
	Amount   int64
}

func (d *NoDoubleTransaction) getLimiter(transaction *model.Transaction) *util.RateLimiter {
	key := doubleTransactionKey{transaction.Merchant, transaction.Amount}
	limiter := d.limiters[key]
	if limiter == nil {
		limiter = util.NewRateLimiter(1, d.doubleTxInterval)
		d.limiters[key] = limiter
	}
	return limiter
}
