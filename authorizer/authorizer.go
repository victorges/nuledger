package authorizer

import (
	"nuledger/authorizer/rules"
	"nuledger/model"
	"nuledger/model/violation"
)

type Authorizer struct {
	accountState *model.Account
	rules        []rules.Rule
}

func NewAuthorizer() *Authorizer {
	return &Authorizer{
		rules: rules.Default(),
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

	a.accountState.AvailableLimit -= transaction.Amount
	for _, commit := range commitFuncs {
		commit()
	}
	return *a.accountState, nil
}
