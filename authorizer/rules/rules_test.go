package rules_test

import (
	"errors"
	"nuledger/authorizer/rules"
	"nuledger/model"
	"nuledger/model/violation"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var startTime = time.Date(2021, time.April, 2, 21, 11, 43, 0, time.Local)

func TestAccountCardActive(t *testing.T) {
	Convey("Given AccountCardActive authorizer function", t, func() {
		Convey("It should authorize accounts with an active card", func() {
			commitFunc, err := rules.AccountCardActive(model.Account{ActiveCard: true}, model.Transaction{})
			So(commitFunc, ShouldBeNil)
			So(err, ShouldBeNil)
		})

		Convey("It should NOT authorize accounts with an inactive card", func() {
			_, err := rules.AccountCardActive(model.Account{ActiveCard: false}, model.Transaction{})
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
			_, err := rules.SufficientLimit(model.Account{AvailableLimit: 50}, model.Transaction{Amount: 100})
			So(err, ShouldNotBeNil)
			So(err, ShouldResemble, violation.ErrorInsufficientLimit)
		})
	})
}

func TestChronologicalOrder(t *testing.T) {
	Convey("Given ChronologicalOrder authorizer", t, func() {
		authzer := &rules.ChronologicalOrder{}

		Convey("It should authorize initial transactions", func() {
			commitFunc, err := authzer.Authorize(model.Account{}, model.Transaction{Time: startTime})
			So(commitFunc, ShouldBeNil)
			So(err, ShouldBeNil)

			Convey("Then it should authorize transactions on the same timestamp", func() {
				commitFunc, err := authzer.Authorize(model.Account{}, model.Transaction{Time: startTime})
				So(commitFunc, ShouldBeNil)
				So(err, ShouldBeNil)
			})
			Convey("It should also authorize transactions on later timestamps", func() {
				commitFunc, err := authzer.Authorize(model.Account{}, model.Transaction{Time: startTime.Add(1 * time.Hour)})
				So(commitFunc, ShouldBeNil)
				So(err, ShouldBeNil)
			})
			Convey("It should NOT authorize transactions on earlier timestamps", func() {
				_, err := authzer.Authorize(model.Account{}, model.Transaction{Time: startTime.Add(-1 * time.Microsecond)})
				So(err, ShouldNotBeNil)
				So(errors.As(err, &violation.Error{}), ShouldBeFalse)
			})
		})
	})
}
