package authorizer

import (
	"errors"
	"fmt"
	"nuledger/iop"
	"nuledger/model"
	"nuledger/model/violation"
)

type Handler struct {
	*Authorizer
}

func NewHandler() iop.DataHandler {
	return &Handler{NewAuthorizer()}
}

type operationType int

const (
	operationTypeUnknown operationType = iota
	operationTypeCreateAccount
	operationTypePerformTransaction
)

func (h *Handler) Handle(op model.OperationInput) (model.StateOutput, error) {
	opType, err := getOperationType(op)
	if err != nil {
		return model.StateOutput{}, fmt.Errorf("Bad operation object: %w", err)
	}

	var account model.Account
	switch opType {
	case operationTypeCreateAccount:
		account, err = h.CreateAccount(op.Account)
	case operationTypePerformTransaction:
		account, err = h.PerformTransaction(op.Transaction)
	default:
		return model.StateOutput{}, errors.New("Internal error: Unknown operation type")
	}

	violations, err := extractViolations(err)
	if err != nil {
		return model.StateOutput{}, err
	}
	return model.StateOutput{account, violations}, nil
}

func getOperationType(op model.OperationInput) (operationType, error) {
	hasAccount := op.Account != nil
	hasTransaction := op.Transaction != nil
	if hasAccount == hasTransaction {
		return operationTypeUnknown, errors.New(`Must have exactly 1 of "account" or "transaction" fields set`)
	}
	if hasAccount {
		return operationTypeCreateAccount, nil
	}
	return operationTypePerformTransaction, nil
}

func extractViolations(err error) ([]violation.Code, error) {
	var verr *violation.Error
	if errors.As(err, &verr) {
		return []violation.Code{verr.Code}, nil
	}
	var aggErr model.AggregateError
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
	return violations, model.AggregateErrors(fatalErrs)
}
