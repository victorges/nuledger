package authorizer_test

import (
	"errors"
	"nuledger/authorizer"
	"nuledger/authorizer/rule"
	mock_rule "nuledger/mocks/authorizer/rule"
	"nuledger/model"
	"nuledger/model/violation"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	ledgerStartTime  = time.Date(2021, time.April, 1, 12, 49, 36, 0, time.Local)
	dummyTransaction = model.Transaction{Merchant: "Ribon App", Amount: 100, Time: ledgerStartTime}
)

func TestLedger(t *testing.T) {
	Convey("Given an authorizer Ledger", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		authzer := mock_rule.NewMockAuthorizer(ctrl)
		ledger := authorizer.NewLedger(authzer)

		Convey("When no account has been created", func() {
			Convey("It should return an error for any transaction perform", func() {
				account, err := ledger.PerformTransaction(dummyTransaction)
				So(err, ShouldNotBeNil)
				So(account, ShouldBeNil)

				var verr violation.Error
				So(errors.As(err, &verr), ShouldBeTrue)
				So(verr.Code, ShouldEqual, violation.AccountNotInitialized)
			})

			Convey("It should allow creating an account", func() {
				accountReq := model.Account{ActiveCard: true, AvailableLimit: 2}

				account, err := ledger.CreateAccount(accountReq)
				So(err, ShouldBeNil)
				So(account, ShouldNotBeNil)
				So(*account, ShouldResemble, accountReq)
			})
		})

		Convey("When there is an account created", func() {
			initAccountState := model.Account{ActiveCard: true, AvailableLimit: 500}
			_, err := ledger.CreateAccount(initAccountState)
			So(err, ShouldBeNil)

			Convey("It should return an error for creating another account", func() {
				accountReq := model.Account{ActiveCard: true, AvailableLimit: 2}
				account, err := ledger.CreateAccount(accountReq)
				So(err, ShouldNotBeNil)
				So(account, ShouldNotBeNil)
				So(*account, ShouldResemble, initAccountState)

				var verr violation.Error
				So(errors.As(err, &verr), ShouldBeTrue)
				So(verr.Code, ShouldEqual, violation.AccountAlreadyInitialized)
			})

			Convey("It should check transactions with authorizer", func() {
				authzer.EXPECT().
					Authorize(gomock.Eq(initAccountState), gomock.Eq(dummyTransaction)).
					Return(nil, nil)
				ledger.PerformTransaction(dummyTransaction)
			})

			Convey("And when authorizer returns success", func() {
				test := func(commit rule.CommitFunc) model.Account {
					authzer.EXPECT().
						Authorize(gomock.Eq(initAccountState), gomock.Eq(dummyTransaction)).
						Return(commit, nil)

					account, err := ledger.PerformTransaction(dummyTransaction)
					So(err, ShouldBeNil)
					So(account, ShouldNotBeNil)
					return *account
				}

				Convey("It should update the account state", func() {
					expectedAfterTx := initAccountState
					expectedAfterTx.AvailableLimit -= dummyTransaction.Amount

					account := test(nil)
					So(account, ShouldResemble, expectedAfterTx)
				})
				Convey("It should call returned commitFunc", func() {
					callCount := 0
					commit := func() { callCount++ }

					test(commit)
					So(callCount, ShouldEqual, 1)
				})
			})

			Convey("When authorizer returns an error", func() {
				returnedErr := errors.New("Custom error")
				test := func(commit rule.CommitFunc) (model.Account, error) {
					authzer.EXPECT().
						Authorize(gomock.Eq(initAccountState), gomock.Eq(dummyTransaction)).
						Return(commit, returnedErr)

					account, err := ledger.PerformTransaction(dummyTransaction)
					So(err, ShouldNotBeNil)
					So(account, ShouldNotBeNil)
					return *account, err
				}

				Convey("It should NOT update the acccount state", func() {
					account, _ := test(nil)
					So(account, ShouldResemble, initAccountState)
				})
				Convey("It should propagate the error", func() {
					_, err := test(nil)
					So(err, ShouldEqual, returnedErr)
				})
				Convey("It should NOT call the commit function", func() {
					callCount := 0
					commit := func() { callCount++ }

					test(commit)
					So(callCount, ShouldEqual, 0)
				})
			})

		})
	})
}
