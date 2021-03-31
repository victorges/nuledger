package main

import (
	"os"

	"nuledger/iop"
	"nuledger/model"
	"nuledger/model/violation"
)

type dummyHandler struct{}

func (d dummyHandler) Handle(op model.OperationInput) (model.StateOutput, error) {
	return model.StateOutput{Account: op.Account, Violations: []violation.Code{}}, nil
}

func main() {
	processor := iop.NewProcessor(os.Stdin, os.Stdout, dummyHandler{})
	if err := processor.Process(); err != nil {
		panic(err)
	}
}
