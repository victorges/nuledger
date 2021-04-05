package rules_test

import (
	"fmt"
	"nuledger/authorizer/rules"
	"nuledger/model"
	"nuledger/model/violation"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var frequencyStartTime = time.Date(2021, time.April, 1, 13, 04, 21, 0, time.Local)

func TestLimitedFrequency(t *testing.T) {
	Convey("Given LimitedFrequency authorizer", t, func() {
		maxTxs := 3
		interval := 1 * time.Minute
		authzer := rules.NewLimitedFrequency(maxTxs, interval)

		Convey("It should authorize initial transactions", func() {
			commitFunc, err := authzer.Authorize(model.Account{}, genTransaction(0))
			So(commitFunc, ShouldNotBeNil)
			So(err, ShouldBeNil)

			testSuccess := func(transaction model.Transaction) {
				commitFunc, err := authzer.Authorize(model.Account{}, transaction)
				So(commitFunc, ShouldNotBeNil)
				So(err, ShouldBeNil)
				commitFunc()
			}

			Convey("Then it should still authorize unlimited transactions if the commitFunc is not called", func() {
				for i := 0; i < 10*maxTxs; i++ {
					authzer.Authorize(model.Account{}, genTransaction(0))
				}
				testSuccess(baseTransacton)
			})
			commitFunc()

			Convey("It SHOULD authorize", func() {
				Convey("Up until the quota in the same timestamp", func() {
					for i := 1; i < maxTxs; i++ {
						testSuccess(genTransaction(0))
					}

					Convey("Then one immediately after the interval", func() {
						testSuccess(genTransaction(interval))
					})
					Convey("Or a little bit after the interval", func() {
						testSuccess(genTransaction(interval + 1))
					})
				})
				Convey("Any amount of transactions with a minimal period", func() {
					period := interval / time.Duration(maxTxs)
					totalDiff := period
					for i := 0; i < 10*maxTxs; i++ {
						testSuccess(genTransaction(totalDiff))
						totalDiff += period
					}
				})
			})

			Convey("And it should NOT authorize", func() {
				testError := func(transaction model.Transaction) {
					commitFunc, err := authzer.Authorize(model.Account{}, transaction)
					So(commitFunc, ShouldBeNil)
					So(err, ShouldNotBeNil)
					So(err, ShouldResemble, violation.ErrorHighFrequencySmallInterval)
				}

				Convey("If the quota is immediately consumed", func() {
					for i := 1; i < maxTxs; i++ {
						testSuccess(genTransaction(0))
					}
					Convey("Right in the same timestamp", func() {
						testError(genTransaction(0))
					})
					Convey("Right before the end of the interval", func() {
						testError(genTransaction(interval - 1))
					})
				})

				Convey("If the quota is consumed right in the end of the interval", func() {
					for i := 1; i < maxTxs; i++ {
						testSuccess(genTransaction(interval - 1))
					}
					testSuccess(genTransaction(interval))

					Convey("In the beginning of the next interval", func() {
						testError(genTransaction(interval))
						testError(genTransaction(interval + 1))
					})
					Convey("A bit before the end of the next interval", func() {
						testError(genTransaction(2*interval - 2))
					})
				})

				Convey("If the quota is consumed a tad too fast", func() {
					period := (interval - 1) / time.Duration(maxTxs)
					totalDiff := period
					for i := 1; i < maxTxs; i++ {
						testSuccess(genTransaction(totalDiff))
						totalDiff += period
					}
					testError(genTransaction(totalDiff))
				})
			})
		})
	})
}

var txIdx int64 = 0

func genTransaction(diff time.Duration) model.Transaction {
	txIdx++
	return model.Transaction{
		Merchant: fmt.Sprintf("Merchant %d", txIdx),
		Amount:   100 + txIdx,
		Time:     frequencyStartTime.Add(diff),
	}
}
