// Package model contains all the model types shared by the whole application.
package model

// Account represents both the current account state sent on response messages
// as well as the account creation object representing its initial state. In a
// multi-account setup, it should include some ID of the account.
type Account struct {
	ActiveCard     bool  `json:"active-card"`
	AvailableLimit int64 `json:"available-limit"`
}

// Copy is a helper function for creating a copy of the current object and
// returning it as a pointer.
func (a Account) Copy() *Account {
	return &a
}
