package authorizer

import (
	"errors"
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

func (h *Handler) Handle(op model.OperationInput) (model.StateOutput, error) {
	if op.Account != nil && op.Transaction != nil {
		return model.StateOutput{}, errors.New("Must have only 1 of account or transaction fields set")
	} else if op.Account != nil {
		account, err := h.CreateAccount(op.Account)
		var verr *violation.Error
		if err != nil && !errors.As(err, &verr) {
			return model.StateOutput{}, err
		}
		return newStateOutput(account, verr), nil
	} else if op.Transaction != nil {
		account, err := h.PerformTransaction(op.Transaction)
		var verr *violation.Error
		if err != nil && !errors.As(err, &verr) {
			return model.StateOutput{}, err
		}
		return newStateOutput(account, verr), nil
	} else {
		return model.StateOutput{}, errors.New("Must have either account or transaction fields set")
	}
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
