package rules

import (
	"nuledger/model"
	"nuledger/model/violation"
)

type AccountCardActive struct{}

func (AccountCardActive) Authorize(account model.Account, transaction *model.Transaction) (CommitFunc, error) {
	if !account.ActiveCard {
		err := violation.NewError(violation.CardNotActive, "Account card is not active")
		return nil, err
	}
	return nil, nil
}

type SufficientLimit struct{}

func (SufficientLimit) Authorize(account model.Account, transaction *model.Transaction) (CommitFunc, error) {
	if account.AvailableLimit < transaction.Amount {
		err := violation.NewError(violation.InsufficientLimit, "Transaction amount is higher than available limit")
		return nil, err
	}
	return nil, nil
}
