package rules

import (
	"nuledger/model"
	"nuledger/model/violation"
	"time"
)

type NoDoubleTransaction struct {
	doubleTxInterval time.Duration
	limiters         map[doubleTransactionKey]*RateLimiter
}

func NewNoDoubleTransaction(interval time.Duration) Rule {
	return &NoDoubleTransaction{interval, map[doubleTransactionKey]*RateLimiter{}}
}

func (d *NoDoubleTransaction) Authorize(_ model.Account, transaction *model.Transaction) (CommitFunc, error) {
	limiter := d.getLimiter(transaction)
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

func (d *NoDoubleTransaction) getLimiter(transaction *model.Transaction) *RateLimiter {
	key := doubleTransactionKey{transaction.Merchant, transaction.Amount}
	limiter := d.limiters[key]
	if limiter == nil {
		limiter = NewRateLimiter(1, d.doubleTxInterval)
		d.limiters[key] = limiter
	}
	return limiter
}
