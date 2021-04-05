package rules_test

import (
	"nuledger/authorizer/rules"
	"nuledger/model"
	"nuledger/model/violation"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	uniqueStartTime = time.Date(2021, time.April, 1, 13, 05, 19, 0, time.Local)
	baseTransacton  = model.Transaction{Merchant: "One Merchant", Amount: 1, Time: uniqueStartTime}
)

func TestUniqueTransactions(t *testing.T) {
	Convey("Given UniqueTransactions authorizer", t, func() {
		interval := 1 * time.Minute
		authzer := rules.NewUniqueTransactions(interval)

		Convey("It should authorize initial transactions", func() {
			commitFunc, err := authzer.Authorize(model.Account{}, baseTransacton)
			So(commitFunc, ShouldNotBeNil)
			So(err, ShouldBeNil)

			testSuccess := func(transaction model.Transaction) {
				commitFunc, err := authzer.Authorize(model.Account{}, transaction)
				So(commitFunc, ShouldNotBeNil)
				So(err, ShouldBeNil)
				commitFunc()
			}

			Convey("Then it should still authorize an identical transaction if the commitFunc is not called", func() {
				testSuccess(baseTransacton)
			})
			commitFunc()

			Convey("It should NOT authorize", func() {
				testError := func(transaction model.Transaction) {
					commitFunc, err := authzer.Authorize(model.Account{}, transaction)
					So(commitFunc, ShouldBeNil)
					So(err, ShouldNotBeNil)
					So(err, ShouldResemble, violation.ErrorDoubleTransaction)
				}

				Convey("Identical transactions", func() {
					testError(baseTransacton)
				})
				Convey("Transactions on too near timestamps", func() {
					repeatedTransaction := baseTransacton
					repeatedTransaction.Time = uniqueStartTime.Add(interval / 2)
					testError(repeatedTransaction)
				})
			})

			Convey("And it SHOULD authorize", func() {
				Convey("Transactions from other merchants", func() {
					otherMerchant := baseTransacton
					otherMerchant.Merchant = "Another Merchant"
					testSuccess(otherMerchant)
				})
				Convey("Transactions of another amount", func() {
					otherAmount := baseTransacton
					otherAmount.Amount = 2 * baseTransacton.Amount
					testSuccess(otherAmount)
				})
				Convey("Transactions after the interval", func() {
					validTransaction := baseTransacton
					validTransaction.Time = uniqueStartTime.Add(interval + 1)
					testSuccess(validTransaction)
				})
			})
		})
	})
}
