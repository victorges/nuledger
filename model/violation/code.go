// Package violation contains useful types for representing the well-defined
// violations in the authorizer business logic.
package violation

// Code represents a violation code to be included in the output messages
// respective to the well-defined errors processing requested operations.
type Code string

const (
	AccountNotInitialized      Code = "account-not-initialized"
	AccountAlreadyInitialized       = "account-already-initialized"
	InsufficientLimit               = "insufficient-limit"
	CardNotActive                   = "card-not-active"
	HighFrequencySmallInterval      = "high-frequency-small-interval"
	DoubleTransaction               = "double-transaction"
)
