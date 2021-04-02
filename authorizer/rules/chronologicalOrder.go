package rules

import (
	"fmt"
	"nuledger/model"
	"time"
)

type ChronologicalOrder struct {
	lastTxTime time.Time
}

func (c *ChronologicalOrder) Validate(_ model.Account, transaction *model.Transaction) (CommitFunc, error) {
	if transaction.Time.Before(c.lastTxTime) {
		return nil, fmt.Errorf("Transactions must be sent in chronological order. Received %v after %v", transaction.Time, c.lastTxTime)
	}
	c.lastTxTime = transaction.Time
	return nil, nil
}
