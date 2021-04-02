package rules

import (
	"fmt"
	"nuledger/authorizer/rule"
	"nuledger/model"
	"time"
)

type ChronologicalOrder struct {
	lastTxTime time.Time
}

func (c *ChronologicalOrder) Authorize(_ model.Account, transaction *model.Transaction) (rule.CommitFunc, error) {
	if transaction.Time.Before(c.lastTxTime) {
		return nil, fmt.Errorf("Transactions must be sent in chronological order. Received %v after %v", transaction.Time, c.lastTxTime)
	}
	c.lastTxTime = transaction.Time
	return nil, nil
}
