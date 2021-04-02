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

func (a *Authorizer) CreateAccount(account *model.Account) (*model.Account, error) {
	if a.accountState != nil {
		return a.accountState.Copy(), violation.ErrorAccountAlreadyInitialized
	}

	a.accountState = account.Copy()
	return account, nil
}

func (a *Authorizer) PerformTransaction(transaction *model.Transaction) (*model.Account, error) {
	if a.accountState == nil {
		return nil, violation.ErrorAccountNotInitialized
	}
	account := a.accountState

	commitFuncs, err := authorize(*account, a.rules, transaction)
	if err != nil {
		return account.Copy(), err
	}

	account.AvailableLimit -= transaction.Amount
	invokeAll(commitFuncs)
	return account.Copy(), nil
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
	return commitFuncs, model.AggregateErrors(errs)
}

func invokeAll(funcs []rules.CommitFunc) {
	for _, f := range funcs {
		f()
	}
}
