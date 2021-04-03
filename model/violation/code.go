// Package violation contains useful types for representing the well-defined
// violations in the authorizer business logic.
package violation

// Code is an enum to represent each of the well-known violation codes that can
// be included in the output messages.
type Code string

const (
	AccountNotInitialized      Code = "account-not-initialized"
	AccountAlreadyInitialized       = "account-already-initialized"
	InsufficientLimit               = "insufficient-limit"
	CardNotActive                   = "card-not-active"
	HighFrequencySmallInterval      = "high-frequency-small-interval"
	DoubleTransaction               = "double-transaction"
)
