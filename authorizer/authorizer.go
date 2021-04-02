package authorizer

import (
	"nuledger/authorizer/rules"
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
	rules            []rules.Rule
}

func NewAuthorizer() *Authorizer {
	return &Authorizer{
		globalLimiter:    NewRateLimiter(3, 2*time.Minute),
		doubleTxLimiters: map[doubleTransactionKey]*RateLimiter{},
		rules:            rules.Default(),
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
	if a.accountState == nil {
		err := violation.NewError(violation.AccountNotInitialized, "Account hasn't been initialized")
		// TODO: Change return to a pointer to have a null output instead of default object
		return model.Account{}, err
	}

	var (
		commitFuncs = make([]rules.CommitFunc, 0, 2)
		errs        []error
	)
	for _, rule := range a.rules {
		commit, err := rule.Validate(*a.accountState, transaction)
		if commit != nil {
			commitFuncs = append(commitFuncs, commit)
		}
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		// TODO: Aggregate errors
		return *a.accountState, errs[0]
	}

	account := a.accountState
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
	for _, commit := range commitFuncs {
		commit()
	}
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
