package util_test

import (
	"errors"
	"nuledger/util"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAggregateError(t *testing.T) {
	Convey("Given an AggregateError", t, func() {
		aggErr := util.AggregateError{
			Errors: []error{
				errors.New("ABC"),
				errors.New("XYZ"),
				errors.New("BTT")},
		}

		Convey("Its Error function should include all messages", func() {
			msg := aggErr.Error()
			for _, err := range aggErr.Errors {
				So(msg, ShouldContainSubstring, err.Error())
			}
		})
	})
}

func TestAggregateErrors(t *testing.T) {
	Convey("Given AggregateErrors function", t, func() {
		Convey("It should return a nil error", func() {
			Convey("For a nil input", func() {
				So(util.AggregateErrors(nil), ShouldBeNil)
			})

			Convey("For an empty slice", func() {
				So(util.AggregateErrors([]error{}), ShouldBeNil)
			})
		})

		Convey("It should the single error if slice of size 1", func() {
			theErr := errors.New("It's me")
			So(util.AggregateErrors([]error{theErr}), ShouldEqual, theErr)
		})

		Convey("It should return an aggregate error for multiple", func() {
			theErr := errors.New("It's me")
			expected := util.AggregateError{Errors: []error{theErr, theErr, theErr}}

			out := util.AggregateErrors([]error{theErr, theErr, theErr})
			So(out, ShouldNotBeNil)
			So(out, ShouldResemble, expected)
		})
	})
}
