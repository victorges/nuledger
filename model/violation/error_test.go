package violation_test

import (
	"nuledger/model/violation"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestError(t *testing.T) {
	Convey("Given an error object, it should", t, func() {
		code := violation.Code("custom-violation-code")
		err := violation.NewError(code, "Custom error message: %v", code)

		Convey("Contain the exact same code within", func() {
			So(err.Code, ShouldEqual, code)
		})

		Convey("Format the message as expected", func() {
			So(err.Message, ShouldEqual, "Custom error message: custom-violation-code")
		})

		Convey("Return the Message field in Error function", func() {
			So(err.Error(), ShouldEqual, err.Message)
		})
	})
}
