package main

import (
	"errors"
	"os"

	"nuledger/authorizer"
	"nuledger/iop"
	"nuledger/model"
	"nuledger/model/violation"
)

type handler struct {
	authorizer.Authorizer
}

func (h handler) Handle(op model.OperationInput) (model.StateOutput, error) {
	if op.Account != nil && op.Transaction != nil {
		return model.StateOutput{}, errors.New("Must have only 1 of account or transaction fields set")
	} else if op.Account != nil {
		account, err := h.CreateAccount(op.Account)
		if err != nil {
			return model.StateOutput{}, err
		}
		return newStateOutput(account), nil
	} else if op.Transaction != nil {
		account, err := h.PerformTransaction(op.Transaction)
		if err != nil {
			return model.StateOutput{}, err
		}
		return newStateOutput(account), nil
	} else {
		return model.StateOutput{}, errors.New("Must have either account or transaction fields set")
	}
}

func newStateOutput(account model.Account) model.StateOutput {
	return model.StateOutput{
		Account:    account,
		Violations: []violation.Code{},
	}
}

func main() {
	processor := iop.NewProcessor(os.Stdin, os.Stdout, handler{authorizer.Authorizer{}})
	if err := processor.Process(); err != nil {
		panic(err)
	}
}
