package authorizer

import (
	"fmt"
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
	lastTxTime       time.Time
}

func NewAuthorizer() *Authorizer {
	return &Authorizer{
		globalLimiter:    NewRateLimiter(3, 2*time.Minute),
		doubleTxLimiters: map[doubleTransactionKey]*RateLimiter{},
	}
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
	if transaction.Time.Before(a.lastTxTime) {
		return model.Account{}, fmt.Errorf("Transactions must be sent in chronological order. Received %v after %v", transaction.Time, a.lastTxTime)
	}
	a.lastTxTime = transaction.Time

	// TODO: Create some kind of authorization rule interface to abstract each of these validations.
	account := a.accountState
	if account == nil {
		err := violation.NewError(violation.AccountNotInitialized, "Account hasn't been initialized")
		// TODO: Change return to a pointer to have a null output instead of default object
		return model.Account{}, err
	}
	if !account.ActiveCard {
		err := violation.NewError(violation.CardNotActive, "Account card is not active")
		return *account, err
	}
	if account.AvailableLimit < transaction.Amount {
		err := violation.NewError(violation.InsufficientLimit, "Transaction amount is higher than available limit")
		return *account, err
	}
	if !a.globalLimiter.Allow(transaction.Time) {
		err := violation.NewError(violation.HighFrequencySmallInterval, "Too many transactions in a small interval")
		return *account, err
	}
	doubleTxLimiter := a.getDoubleTransactionLimiter(transaction)
	if !doubleTxLimiter.Allow(transaction.Time) {
		err := violation.NewError(violation.DoubleTransaction, "Duplicate transaction of same amount and merchant")
		return *account, err
	}
	a.globalLimiter.Take(transaction.Time)
	doubleTxLimiter.Take(transaction.Time)
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
