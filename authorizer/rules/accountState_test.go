package rules_test

import (
	"nuledger/authorizer/rules"
	"nuledger/model"
	"nuledger/model/violation"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var startTime = time.Date(2021, time.March, 31, 21, 11, 43, 0, time.Local)

func TestAccountCardActive(t *testing.T) {
	Convey("Given AccountCardActive authorizer function", t, func() {
		Convey("It should authorize accounts with an active card", func() {
			commitFunc, err := rules.AccountCardActive(model.Account{ActiveCard: true}, model.Transaction{})
			So(commitFunc, ShouldBeNil)
			So(err, ShouldBeNil)
		})

		Convey("It should NOT authorize accounts with an inactive card", func() {
			commitFunc, err := rules.AccountCardActive(model.Account{ActiveCard: false}, model.Transaction{})
			So(commitFunc, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err, ShouldResemble, violation.ErrorCardNotActive)
		})
	})
}

func TestSufficientLimit(t *testing.T) {
	Convey("Given SufficientLimit authorizer function", t, func() {
		Convey("It should authorize accounts with sufficient limit", func() {
			commitFunc, err := rules.SufficientLimit(model.Account{AvailableLimit: 50}, model.Transaction{Amount: 10})
			So(commitFunc, ShouldBeNil)
			So(err, ShouldBeNil)
		})

		Convey("It should NOT authorize accounts with an insufficient limit", func() {
			commitFunc, err := rules.SufficientLimit(model.Account{AvailableLimit: 50}, model.Transaction{Amount: 100})
			So(commitFunc, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err, ShouldResemble, violation.ErrorInsufficientLimit)
		})
	})
}
