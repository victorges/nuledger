package rules

import (
	"nuledger/authorizer/rule"
	"nuledger/model"
	"nuledger/model/violation"
	"nuledger/util"
	"time"
)

// UniqueTransactions is a rule.Authorizer to guarantee that no double/duplicate
// transactions are allowed to go through in the account. For 2 transactions to
// be considered double, they must have the exact same merchant and amount, and
// have timestamps within the configured interval of each other.
type UniqueTransactions struct {
	doubleTxInterval time.Duration
	limiters         map[doubleTransactionKey]*util.RateLimiter
}

// NewUniqueTransactions returns a UniqueTransactions authorizer with the given
// configuration. The `interval` is the minimum amount of time between two
// transactions with everything else equal for them not to glbe considered double.
func NewUniqueTransactions(interval time.Duration) *UniqueTransactions {
	return &UniqueTransactions{interval, map[doubleTransactionKey]*util.RateLimiter{}}
}

// Authorize checks if the given transaction is a double, considering the
// configured interval, and if so the transaction is not authorized and a
// double-transaction violation error is returned.
func (d *UniqueTransactions) Authorize(_ model.Account, transaction model.Transaction) (rule.CommitFunc, error) {
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

// getLimiter tries to get the existing rate limiter for a given transaction and
// creates a new one if there is none yet.
func (d *UniqueTransactions) getLimiter(transaction *model.Transaction) *util.RateLimiter {
	key := doubleTransactionKey{transaction.Merchant, transaction.Amount}
	limiter := d.limiters[key]
	if limiter == nil {
		limiter = &util.RateLimiter{MaxEvents: 1, Interval: d.doubleTxInterval}
		d.limiters[key] = limiter
	}
	return limiter
}
