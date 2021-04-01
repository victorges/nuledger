package authorizer

import (
	"nuledger/model"
	"nuledger/model/violation"
)

type Authorizer struct {
	accountState   *model.Account
	globalLimiter  *RateLimiter
}

func (a *Authorizer) CreateAccount(account *model.Account) (model.Account, error) {
	if a.accountState != nil {
		err := violation.NewError(violation.AccountAlreadyInitialized, "Account has already been initialized")
		return *a.accountState, err
	}

	a.accountState = &model.Account{}
	*a.accountState = *account
	return *a.accountState, nil
}

func (a *Authorizer) PerformTransaction(transaction *model.Transaction) (model.Account, error) {
	account := a.accountState
	if account == nil {
		err := violation.NewError(violation.AccountNotInitialized, "Account hasn't been initialized")
		return model.Account{}, err
	}
	if !account.ActiveCard {
		err := violation.NewError(violation.CardNotActive, "Account card is not active")
		return model.Account{}, err
	}
	if account.AvailableLimit < transaction.Amount {
		err := violation.NewError(violation.InsufficientLimit, "Transaction amount is higher than available limit")
		return model.Account{}, err
	}
	if ok, err := a.globalLimiter.Take(transaction.Time); err != nil {
		return model.Account{}, err
	} else if !ok {
		err := violation.NewError(violation.HighFrequencySmallInterval, "Too many transactions in a small interval")
		return model.Account{}, err
	}
	account.AvailableLimit -= transaction.Amount
	return *account, nil
}
