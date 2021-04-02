package iop

import (
	"nuledger/model"
	"nuledger/model/violation"
)

// OperationInput is a JSON received as an input for an operation to be run.
type OperationInput struct {
	Account     *model.Account     `json:"account"`
	Transaction *model.Transaction `json:"transaction"`
}

// StateOutput represents a JSON to be written in the output as the result of
// performing an operation.
type StateOutput struct {
	Account    *model.Account   `json:"account"`
	Violations []violation.Code `json:"violations"`
}