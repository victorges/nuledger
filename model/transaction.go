package model

import "time"

// Transaction is an authorization request for a transaction, consisting of
// information about the merchant, amount to be charged and time. In a multi-
// account setup it should include some ID of the account doing the transaction.
type Transaction struct {
	// Merchant is a unique string to represent the merchant with which a
	// transaction is being made.
	Merchant string `json:"merchant"`
	// Amount is the units of currency that the transaction would be consuming.
	Amount int64 `json:"amount"`
	// Time is the exact time on which the transaction was attempted.
	Time time.Time `json:"time"`
}
