package model_test

import (
	"nuledger/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAccountCopy(t *testing.T) {
	Convey("Given an account object", t, func() {
		account := &model.Account{ActiveCard: true, AvailableLimit: 42439}

		Convey("Copy creates a new object", func() {
			copy := account.Copy()
			So(copy, ShouldNotEqual, account)

			Convey("With all the same fields", func() {
				So(*copy, ShouldResemble, *account)
			})
		})
	})
}
