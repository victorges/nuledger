package rules

import (
	"nuledger/model"
	"nuledger/model/violation"
)

type AccountCardActive struct{}

func (AccountCardActive) Authorize(account model.Account, transaction *model.Transaction) (CommitFunc, error) {
	if !account.ActiveCard {
		return nil, violation.ErrorCardNotActive
	}
	return nil, nil
}

type SufficientLimit struct{}

func (SufficientLimit) Authorize(account model.Account, transaction *model.Transaction) (CommitFunc, error) {
	if account.AvailableLimit < transaction.Amount {
		return nil, violation.ErrorInsufficientLimit
	}
	return nil, nil
}
