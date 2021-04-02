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

	commitFuncs, err := authorize(*a.accountState, a.rules, transaction)
	if err != nil {
		return *a.accountState, err
	}

	a.accountState.AvailableLimit -= transaction.Amount
	invokeAll(commitFuncs)
	return *a.accountState, nil
}

func authorize(account model.Account, authRules []rules.Rule, transaction *model.Transaction) ([]rules.CommitFunc, error) {
	var (
		commitFuncs = make([]rules.CommitFunc, 0, 2)
		errs        []error
	)
	for _, rule := range authRules {
		commit, err := rule.Authorize(account, transaction)
		if commit != nil {
			commitFuncs = append(commitFuncs, commit)
		}
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		// TODO: Aggregate errors
		return commitFuncs, errs[0]
	}
	return commitFuncs, nil
}

func invokeAll(funcs []rules.CommitFunc) {
	for _, f := range funcs {
		f()
	}
}
