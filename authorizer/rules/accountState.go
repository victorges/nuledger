package rules

import (
	"nuledger/model"
	"nuledger/model/violation"
)

func AccountCardActive(account model.Account, transaction *model.Transaction) (CommitFunc, error) {
	if !account.ActiveCard {
		return nil, violation.ErrorCardNotActive
	}
	return nil, nil
}

func SufficientLimit(account model.Account, transaction *model.Transaction) (CommitFunc, error) {
	if account.AvailableLimit < transaction.Amount {
		return nil, violation.ErrorInsufficientLimit
	}
	return nil, nil
}
