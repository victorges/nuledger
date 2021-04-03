package rules

import (
	"fmt"
	"nuledger/authorizer/rule"
	"nuledger/model"
	"time"
)

// ChronologicalOrder is a rule.Authorizer to enforce that all transactions are
// received in chronological order, i.e. with ascending timestamps. It is a
// sanity check of an inherent assumption of the whole system, since many of the
// algorithms won't work if the transactions are not ordered by timestamps.
type ChronologicalOrder struct {
	lastTxTime time.Time
}

// Authorize checks if the current transaction has a timestamp greater than the
// last received transation, and returns a regular (fatal) error if it does not.
func (c *ChronologicalOrder) Authorize(_ model.Account, transaction model.Transaction) (rule.CommitFunc, error) {
	if transaction.Time.Before(c.lastTxTime) {
		return nil, fmt.Errorf("Transactions must be sent in chronological order. Received %v after %v", transaction.Time, c.lastTxTime)
	}
	c.lastTxTime = transaction.Time
	return nil, nil
}
