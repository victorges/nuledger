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

// NewUniqueTransactions returns a rule.Authorizer to guarantee that no double
// (duplicate) transactions are allowed to go through in the account. For 2
// transactions to be considered double, they must have the exact same merchant
// and amount, and have timestamps within the configured interval of each other.
//
// The `interval` parameter is the minimum amount of time between two
// transactions with everything else equal for them not to be considered double.
//
// Its Authorize function checks if the given transaction is a double,
// considering the configured interval, and if so the transaction is not
// authorized and a double-transaction violation error is returned.
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
