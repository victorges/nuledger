package model

import "time"

// Transaction is an authorization request for a transaction, consisting of
// information about the merchant, amount to be charged and time.
type Transaction struct {
	// AccountID is the unique identifier of the account perforing the
	// respective transaction.
	AccountID string `json:"accountId"`
	// Merchant is a unique string to represent the merchant with which a
	// transaction is being made.
	Merchant string `json:"merchant"`
	// Amount is the units of currency that the transaction would be consuming.
	Amount int64 `json:"amount"`
	// Time is the exact time on which the transaction was attempted.
	Time time.Time `json:"time"`
}
