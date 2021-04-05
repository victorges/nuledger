package authorizer_test

import (
	"errors"
	"fmt"
	"nuledger/authorizer"
	"nuledger/authorizer/rule"
	"nuledger/authorizer/rules"
	"nuledger/iop"
	mock_authorizer "nuledger/mocks/authorizer"
	"nuledger/model"
	"nuledger/model/violation"
	"nuledger/util"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

var startTime = time.Date(2021, time.March, 31, 21, 11, 43, 0, time.Local)

func TestHandler(t *testing.T) {
	Convey("Given an authorizer Handler", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		Convey("When it is sent an acceptable input", func() {
			Convey("And ledger returns a successful response", func() {
				returnedAccount := &model.Account{ActiveCard: true, AvailableLimit: 20170415}
				expected := iop.StateOutput{Account: returnedAccount, Violations: []violation.Code{}}

				Convey("It should forward exact response in output", func() {
					validate := func(output iop.StateOutput, err error) {
						So(err, ShouldBeNil)
						So(output, ShouldResemble, expected)
					}
					testHandlerOperations(ctrl, validate, returnedAccount, nil)
				})
			})

			Convey("And ledger returns an error", func() {
				Convey("It should return any violation.Error as a violation", func() {
					violErr := violation.NewError("custom-validation-code", "Hello violations")
					expected := iop.StateOutput{Account: nil, Violations: []violation.Code{violErr.Code}}

					validate := func(output iop.StateOutput, err error) {
						So(err, ShouldBeNil)
						So(output, ShouldResemble, expected)
					}
					testHandlerOperations(ctrl, validate, nil, violErr)
				})

				Convey("It should handle aggregated violation errors", func() {
					returnedError := util.AggregateError{[]error{
						violation.NewError("custom-validation-code", "Hello violations"),
						violation.NewError("yet-another-validation-code", "Old friend"),
					}}
					expected := iop.StateOutput{Account: nil, Violations: []violation.Code{"custom-validation-code", "yet-another-validation-code"}}

					validate := func(output iop.StateOutput, err error) {
						So(err, ShouldBeNil)
						So(output, ShouldResemble, expected)
					}
					testHandlerOperations(ctrl, validate, nil, returnedError)
				})

				Convey("Any other error should be propagated", func() {
					regularErr := errors.New("This is just a regular error")

					validate := func(output iop.StateOutput, err error) {
						So(err, ShouldNotBeNil)
						So(output, ShouldBeZeroValue)
						So(err, ShouldEqual, regularErr)
					}
					testHandlerOperations(ctrl, validate, nil, regularErr)
				})
				Convey("Even if aggregated with other violation errors", func() {
					regularErr := errors.New("This is just a regular error")
					returnedError := util.AggregateError{[]error{
						violation.NewError("custom-validation-code", "Hello violations"),
						regularErr,
					}}

					validate := func(output iop.StateOutput, err error) {
						So(err, ShouldNotBeNil)
						So(output, ShouldBeZeroValue)
						So(err, ShouldEqual, regularErr)
					}
					testHandlerOperations(ctrl, validate, nil, returnedError)
				})
			})
		})
	})
}

func TestHandlerBadInput(t *testing.T) {
	Convey("Given the authorizer Handler gets some bad input", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ledger := mock_authorizer.NewMockLedger(ctrl)
		var handler iop.DataHandler = &authorizer.Handler{ledger}

		Convey("It should return a fatal error", func() {
			test := func(input iop.OperationInput) {
				output, err := handler.Handle(input)
				So(err, ShouldNotBeNil)
				So(output, ShouldBeZeroValue)
			}

			Convey("For an empty input", func() {
				test(iop.OperationInput{})
			})
			Convey("For ambiguous operations", func() {
				test(iop.OperationInput{&model.Account{}, &model.Transaction{}})
			})
		})

		Convey("It should still forward the request", func() {
			uniqueAccount := &model.Account{ActiveCard: true, AvailableLimit: 812674182736172}
			expected := iop.StateOutput{Account: uniqueAccount, Violations: []violation.Code{}}

			test := func(input iop.OperationInput) {
				output, err := handler.Handle(input)
				So(err, ShouldBeNil)
				So(output, ShouldResemble, expected)
			}
			testBothOperations := func(account *model.Account, transaction *model.Transaction) {
				Convey("For CreateAccount operation", func() {
					ledger.EXPECT().
						CreateAccount(gomock.Eq(*account)).
						Return(uniqueAccount, nil)
					test(iop.OperationInput{Account: account})
				})
				Convey("For PerformTransaction operation", func() {
					ledger.EXPECT().
						PerformTransaction(gomock.Eq(*transaction)).
						Return(uniqueAccount, nil)
					test(iop.OperationInput{Transaction: transaction})
				})
			}

			Convey("For empty operation objects", func() {
				testBothOperations(&model.Account{}, &model.Transaction{})
			})
			Convey("For semantically bad operation objects", func() {
				badAccount := &model.Account{AvailableLimit: -23}
				badTransaction := &model.Transaction{Amount: -42}
				testBothOperations(badAccount, badTransaction)
			})
		})
	})
}

func TestDefaultAuthorizers(t *testing.T) {
	Convey("Given the default authorizers", t, func() {
		authzer := authorizer.DefaultAuthorizer()

		Convey("They should be a rule list", func() {
			So(authzer, ShouldHaveSameTypeAs, rule.List{})

			list := authzer.(rule.List)

			Convey("With all required authorization rules", func() {
				So(list, ShouldHaveLength, 5)
				So(list, ShouldContain, &rules.ChronologicalOrder{})
				So(list, ShouldContain, rules.NewLimitedFrequency(3, 2*time.Minute))
				So(list, ShouldContain, rules.NewUniqueTransactions(2*time.Minute))
				So(containsAuthFunc(list, rules.AccountCardActive), ShouldBeTrue)
				So(containsAuthFunc(list, rules.SufficientLimit), ShouldBeTrue)
			})
		})
	})

	Convey("Given a default handler", t, func() {
		handler := authorizer.NewHandler()
		_, err := handler.Handle(iop.OperationInput{Account: &model.Account{}})
		So(err, ShouldBeNil)

		Convey("It should use the default authorizers", func() {
			expected := iop.StateOutput{Account: &model.Account{}, Violations: []violation.Code{violation.CardNotActive}}
			output, err := handler.Handle(iop.OperationInput{Transaction: &model.Transaction{}})
			So(err, ShouldBeNil)
			So(output, ShouldResemble, expected)
		})
	})
}

func containsAuthFunc(list rule.List, authFunc rule.AuthorizerFunc) bool {
	for _, auth := range list {
		if elmFunc, isFunc := auth.(rule.AuthorizerFunc); isFunc {
			if fmt.Sprint(authFunc) == fmt.Sprint(elmFunc) {
				return true
			}
		}
	}
	return false
}

func testHandlerOperations(ctrl *gomock.Controller, validate func(iop.StateOutput, error), returnAccount *model.Account, returnErr error) {
	ledger := mock_authorizer.NewMockLedger(ctrl)
	var handler iop.DataHandler = &authorizer.Handler{ledger}

	Convey("For CreateAccount (Account) operation", func() {
		account := &model.Account{ActiveCard: true, AvailableLimit: 20210902}
		createAccountOp := iop.OperationInput{Account: account}

		ledger.EXPECT().
			CreateAccount(gomock.Eq(*account)).
			Return(returnAccount, returnErr)

		validate(handler.Handle(createAccountOp))
	})
	Convey("For PerformTransaction (Transaction) operation", func() {
		transaction := &model.Transaction{Merchant: "Amazon Web Services", Amount: 142, Time: startTime}
		performTxOp := iop.OperationInput{Transaction: transaction}

		ledger.EXPECT().
			PerformTransaction(gomock.Eq(*transaction)).
			Return(returnAccount, returnErr)

		validate(handler.Handle(performTxOp))
	})
}
