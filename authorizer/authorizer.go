package authorizer

import (
	"nuledger/model"
	"nuledger/model/violation"
	"time"
)

type doubleTransactionKey struct {
	Merchant string
	Amount   int
}

type Authorizer struct {
	accountState     *model.Account
	globalLimiter    *RateLimiter
	doubleTxLimiters map[doubleTransactionKey]*RateLimiter
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
	doubleTxLimiter := a.getDoubleTransactionLimiter(transaction)
	if ok, err := doubleTxLimiter.Take(transaction.Time); err != nil {
		return model.Account{}, err
	} else if !ok {
		err := violation.NewError(violation.DoubleTransaction, "Duplicate transaction of same amount and merchant")
		return model.Account{}, err
	}
	account.AvailableLimit -= transaction.Amount
	return *account, nil
}

func (a *Authorizer) getDoubleTransactionLimiter(transaction *model.Transaction) *RateLimiter {
	key := doubleTransactionKey{transaction.Merchant, transaction.Amount}
	doubleTxLimiter := a.doubleTxLimiters[key]
	if doubleTxLimiter == nil {
		doubleTxLimiter = NewRateLimiter(1, 2*time.Minute)
		a.doubleTxLimiters[key] = doubleTxLimiter
	}
	return doubleTxLimiter
}
