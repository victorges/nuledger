package model

import "time"

// Transaction is an authorization request for a transaction, consisting of
// information about the merchant, amount to be charged and time. In a multi-
// account setup it should include some ID of the account doing the transaction.
type Transaction struct {
	Merchant string    `json:"merchant"`
	Amount   int64     `json:"amount"`
	Time     time.Time `json:"time"`
}
