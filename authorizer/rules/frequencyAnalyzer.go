package rules

import (
	"nuledger/authorizer/rule"
	"nuledger/model"
	"nuledger/model/violation"
	"nuledger/util"
)

// FrequencyAnalyzer is a generic rule.Authorizer that can be used for any kind
// of authorization rule based on the frequency of transactions given a certain
// constraint.
//
// This constraint on the transactions allow to limit multiple sets of
// transactions in an independent manner, so that each those interfere with the
// frequency limits of the other. That constraint is implemented via a
// key-mapper function which should receive a transaction and return a unique
// value (that can be used as a map key) representing the group the transaction
// should be included in. Examples of this are analyising each account
// transactions or each merchant's transactions separately.
//
// The frequency of transactions is limited via a util.RateLimiter, for which a
// base limiter should be provided as a template for any internal rate limiters
// that may need to be created for new transaction groups.
type FrequencyAnalyzer struct {
	baseLimiter *util.RateLimiter
	keyMapper   func(*model.Transaction) interface{}
	limiters    map[interface{}]*util.RateLimiter
	violation   violation.Error
}

// NewFrequencyAnalyzer creates a new frequency analyzer authorizer which limits
// the frequency of received transactions within their corresponding group. The
// frequency is configured via the `baseLimiter` provided, which is copied when
// a new transaction group is created. The grouping of transactions is made by
// the `keyMapper` function which should return a unique (map-key) value for
// each transaction. Finally, when the rate is exceeded its Authorizer function
// returns the error provided as the last `violation` argument.
func NewFrequencyAnalyzer(baseLimiter util.RateLimiter, keyMapper func(*model.Transaction) interface{}, violation violation.Error) *FrequencyAnalyzer {
	return &FrequencyAnalyzer{
		baseLimiter: &baseLimiter,
		keyMapper:   keyMapper,
		limiters:    map[interface{}]*util.RateLimiter{},
		violation:   violation,
	}
}

// Authorize checks if the given transaction is a exceeds the limit of its
// corresponding transaction group, and if so the transaction is not authorized
// and the violation error configured for this analyzer is returned.
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
