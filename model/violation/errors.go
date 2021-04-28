package violation

var (
	ErrorAccountAlreadyInitialized  = NewError(AccountAlreadyInitialized, "Account has already been initialized")
	ErrorAccountNotInitialized      = NewError(AccountNotInitialized, "Account hasn't been initialized")
	ErrorCardNotActive              = NewError(CardNotActive, "Account card is not active")
	ErrorInsufficientLimit          = NewError(InsufficientLimit, "Transaction amount is higher than available limit")
	ErrorHighFrequencySmallInterval = NewError(HighFrequencySmallInterval, "Too many transactions in a small interval")
	ErrorDoubleTransaction          = NewError(DoubleTransaction, "Duplicate transaction of same amount and merchant")
	ErrorMerchantDenied             = NewError(MerchantDenied, "Merchant is denied any transaction")
)
