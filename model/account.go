package model

// Account represents both the current account state sent on response messages
// as well as the account creation object representing its initial state. In a
// multi-account setup, it should include some ID of the account.
type Account struct {
	ActiveCard     bool `json:"active-card"`
	AvailableLimit int  `json:"available-limit"`
}
