package rules

import (
	"nuledger/authorizer/rule"
	"nuledger/model"
	"nuledger/model/violation"
	"nuledger/util"
	"time"
)

type doubleTransactionKey struct {
	AccountID string
	Merchant  string
	Amount    int64
}

// NewUniqueTransactions returns a UniqueTransactions authorizer with the given
// configuration. The `interval` is the minimum amount of time between two
// transactions with everything else equal for them not to glbe considered double.
func NewUniqueTransactions(interval time.Duration) rule.Authorizer {
	limiter := util.RateLimiter{MaxEvents: 1, Interval: interval}
	keyMapper := func(tx *model.Transaction) interface{} {
		return doubleTransactionKey{
			AccountID: tx.AccountID,
			Merchant:  tx.Merchant,
			Amount:    tx.Amount,
		}
	}
	return NewFrequencyAnalyzer(limiter, keyMapper, violation.ErrorDoubleTransaction)
}
