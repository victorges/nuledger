// Package model contains all the model types shared by the whole application.
package model

// Account represents both the current account state sent on response messages
// as well as the account creation object representing its initial state. In a
// multi-account setup, it should include some ID of the account.
type Account struct {
	// ActiveCard represents if the account card is active or not. An inactive
	// card does not authorize any transactions.
	ActiveCard bool `json:"active-card"`
	// AvailableLimit is the units of currency that the account still has.
	// Transactions consume from this limit and it can never be exceeded.
	AvailableLimit int64    `json"available-limit"`
	DenyList       []string `json:"deny-list"`
}

// Copy is a helper function for creating a copy of the current object and
// returning it as a pointer.
func (a Account) Copy() *Account {
	return &a
}
