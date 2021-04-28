package rules

import (
	"nuledger/authorizer/rule"
	"nuledger/model"
	"nuledger/model/violation"
)

func MerchantDenyList(account model.Account, transaction model.Transaction) (rule.CommitFunc, error) {
	for _, deniedMerchant := range account.DenyList {
		if deniedMerchant == transaction.Merchant {
			return nil, violation.ErrorMerchantDenied
		}
	}
	return nil, nil
}
