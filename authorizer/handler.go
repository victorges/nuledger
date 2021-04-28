package authorizer

import (
	"errors"
	"fmt"
	"nuledger/iop"
	"nuledger/model"
	"nuledger/model/violation"
	"nuledger/util"
)

// Handler is a pipe between the actual raw objects returned and received by the
// input/output processing pipeline and the actual ledger system which manages
// the account state and performs transactions. It basically interprets the JSON
// objects received and calls the correct higher-level APIs from the Ledger.
type Handler struct {
	Ledger
}

// NewHandler creates a new Handler with a Ledger with all the default
// authorizers from DefaultAuthorizer.
func NewHandler() iop.DataHandler {
	return &Handler{NewLedger(DefaultAuthorizer())}
}

// Handle implements the iop.DataHandler interface, receiving JSON objects
// parsed from the input stream, detecting the type of operation that should be
// performed from them, calling the correct API in the Ledger and returning the
// output, also translated to the output that the I/O pipeline expects.
func (h *Handler) Handle(op iop.OperationInput) (iop.StateOutput, error) {
	opType, err := getOperationType(op)
	if err != nil {
		return iop.StateOutput{}, fmt.Errorf("Bad operation object: %w", err)
	}

	var account *model.Account
	switch opType {
	case operationTypeCreateAccount:
		account, err = h.CreateAccount(*op.Account)
	case operationTypePerformTransaction:
		account, err = h.PerformTransaction(*op.Transaction)
	case operationTypeDenyList:
		account, err = h.SetDenyList(op.DenyList)
	}

	violations, err := extractViolations(err)
	if err != nil {
		return iop.StateOutput{}, err
	}
	return iop.StateOutput{account, violations}, nil
}

// operationType is a helper enum to identify the kind of operation to be
// performed (too bad there's no pattern-matching like construct in Go).
type operationType int

const (
	operationTypeUnknown operationType = iota
	operationTypeCreateAccount
	operationTypePerformTransaction
	operationTypeDenyList
)

// getOperationType receives the input JSON object and returns what is the
// requested operation that should be performed from it. It also returns errors
// in case of any semantic issues with the object (e.g. specifying multiple
// operations or none of them).
func getOperationType(op iop.OperationInput) (operationType, error) {
	hasAccount := op.Account != nil
	hasTransaction := op.Transaction != nil
	hasDenyList := op.DenyList != nil
	if countTrues(hasAccount, hasTransaction, hasDenyList) != 1 {
		return operationTypeUnknown, errors.New(`Must have exactly 1 of "account" or "transaction" fields set`)
	}
	if hasAccount {
		return operationTypeCreateAccount, nil
	} else if hasDenyList {
		return operationTypeDenyList, nil
	} else {
		return operationTypePerformTransaction, nil
	}
}

func countTrues(booleans ...bool) int {
	count := 0
	for _, b := range booleans {
		if b {
			count++
		}
	}
	return count
}

// extractViolations receives an error and tries to fetch the specific violation
// codes that might be represented by it. If there are any errors that are not
// violation codes they are returned in the second return value.
func extractViolations(err error) ([]violation.Code, error) {
	var verr violation.Error
	if errors.As(err, &verr) {
		return []violation.Code{verr.Code}, nil
	}
	var aggErr util.AggregateError
	if !errors.As(err, &aggErr) {
		return []violation.Code{}, err
	}

	var (
		violations = make([]violation.Code, 0, len(aggErr.Errors))
		fatalErrs  []error
	)
	for _, innerErr := range aggErr.Errors {
		if errors.As(innerErr, &verr) {
			violations = append(violations, verr.Code)
		} else {
			fatalErrs = append(fatalErrs, innerErr)
		}
	}
	return violations, util.AggregateErrors(fatalErrs)
}
