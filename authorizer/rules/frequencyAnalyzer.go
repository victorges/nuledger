package rules

import (
	"nuledger/authorizer/rule"
	"nuledger/model"
	"nuledger/model/violation"
	"nuledger/util"
)

// UniqueTransactions is a rule.Authorizer to guarantee that no double/duplicate
// transactions are allowed to go through in the account. For 2 transactions to
// be considered double, they must have the exact same merchant and amount, and
// have timestamps within the configured interval of each other.
type FrequencyAnalyzer struct {
	baseLimiter *util.RateLimiter
	keyMapper   func(*model.Transaction) interface{}
	limiters    map[interface{}]*util.RateLimiter
	violation   violation.Error
}

func NewFrequencyAnalyzer(baseLimiter util.RateLimiter, keyMapper func(*model.Transaction) interface{}, violation violation.Error) *FrequencyAnalyzer {
	return &FrequencyAnalyzer{
		baseLimiter: &baseLimiter,
		keyMapper:   keyMapper,
		limiters:    map[interface{}]*util.RateLimiter{},
		violation:   violation,
	}
}

// Authorize checks if the given transaction is a double, considering the
// configured interval, and if so the transaction is not authorized and a
// double-transaction violation error is returned.
func (d *FrequencyAnalyzer) Authorize(_ model.Account, transaction model.Transaction) (rule.CommitFunc, error) {
	limiter := d.getLimiter(&transaction)
	if !limiter.Allows(transaction.Time) {
		return nil, d.violation
	}
	commit := func() { limiter.Take(transaction.Time) }
	return commit, nil
}

// getLimiter tries to get the existing rate limiter for a given transaction and
// creates a new one if there is none yet.
func (d *FrequencyAnalyzer) getLimiter(transaction *model.Transaction) *util.RateLimiter {
	key := d.keyMapper(transaction)
	limiter := d.limiters[key]
	if limiter != nil {
		return limiter
	}

	copy := *d.baseLimiter
	limiter = &copy
	d.limiters[key] = limiter
	return limiter
}
