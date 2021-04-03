package iop

import (
	"nuledger/model"
	"nuledger/model/violation"
)

// OperationInput is a JSON received as an input for an operation to be run. It
// can have only one of Account or Transaction fields set, each one representing
// a different kind of operation being requested.
type OperationInput struct {
	// Account represents an account creation request. If it is not null, it
	// should contain the initial state of the account to be created.
	Account *model.Account `json:"account"`
	// Transaction represents a transaction request. If it is not null, it
	// should contain the details about the transaction being attempted.
	Transaction *model.Transaction `json:"transaction"`
}

// StateOutput represents a JSON to be written in the output as the result of
// performing an operation.
type StateOutput struct {
	// Account represents the current state of the account corresponding to the
	// requested operation. If the operation succeeded it will be the state of
	// the account after the operation, otherwise it will be the current state
	// of the account that couldn't be updated due to the violations.
	Account *model.Account `json:"account"`
	// Violations represent any violation that may have prevented the operation
	// from being performed. It will be an empty array in case of a success.
	Violations []violation.Code `json:"violations"`
}
