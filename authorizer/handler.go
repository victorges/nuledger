package authorizer

import (
	"errors"
	"fmt"
	"nuledger/iop"
	"nuledger/model"
	"nuledger/model/violation"
)

type Handler struct {
	Authorizer
}

func NewHandler() iop.DataHandler {
	return &Handler{Authorizer{}}
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

	var verr *violation.Error
	if err != nil && !errors.As(err, &verr) {
		return model.StateOutput{}, err
	}
	return newStateOutput(account, verr), nil
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

func newStateOutput(account model.Account, err *violation.Error) model.StateOutput {
	out := model.StateOutput{
		Account:    account,
		Violations: []violation.Code{},
	}
	if err != nil {
		out.Violations = append(out.Violations, err.Code)
	}
	return out
}
