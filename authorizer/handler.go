package authorizer

import (
	"errors"
	"fmt"
	"nuledger/iop"
	"nuledger/model"
	"nuledger/model/violation"
	"nuledger/util"
)

type Handler struct {
	*Ledger
}

func NewHandler() iop.DataHandler {
	return &Handler{NewLedger(DefaultAuthorizer())}
}

type operationType int

const (
	operationTypeUnknown operationType = iota
	operationTypeCreateAccount
	operationTypePerformTransaction
)

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
	default:
		return iop.StateOutput{}, errors.New("Internal error: Unknown operation type")
	}

	violations, err := extractViolations(err)
	if err != nil {
		return iop.StateOutput{}, err
	}
	return iop.StateOutput{account, violations}, nil
}

func getOperationType(op iop.OperationInput) (operationType, error) {
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
