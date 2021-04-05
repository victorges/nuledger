package rule_test

import (
	"errors"
	"nuledger/authorizer/rule"
	mock_rule "nuledger/mocks/authorizer/rule"
	"nuledger/model"
	"nuledger/util"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	startTime        = time.Date(2021, time.April, 1, 16, 20, 0, 0, time.Local)
	dummyAccount     = model.Account{ActiveCard: true, AvailableLimit: 234}
	dummyTransaction = model.Transaction{Merchant: "Sketchy", Amount: 123, Time: startTime}
)

func TestRuleList(t *testing.T) {
	Convey("Given a rule List", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		authMocks := make([]*mock_rule.MockAuthorizer, 5)
		authzers := make([]rule.Authorizer, 5)
		for i := range authzers {
			authMocks[i] = mock_rule.NewMockAuthorizer(ctrl)
			authzers[i] = authMocks[i]
		}
		list := rule.List(authzers)

		Convey("It should call all authorizers", func() {
			configureMocks(authMocks)
			_, err := list.Authorize(dummyAccount, dummyTransaction)
			So(err, ShouldBeNil)
		})

		Convey("When there are errors", func() {
			Convey("It should propagate any single error", func() {
				returnedErr := errors.New("Custom error")

				configureMocks(authMocks, 3)
				configureMocksToErr(returnedErr, authMocks[3])

				_, err := list.Authorize(dummyAccount, dummyTransaction)
				So(err, ShouldEqual, returnedErr)
			})

			Convey("It should aggregate multiple errors", func() {
				returnedErr := errors.New("Custom error")
				expectedErr := util.AggregateError{[]error{returnedErr, returnedErr}}

				configureMocks(authMocks, 2, 3)
				configureMocksToErr(returnedErr, authMocks[2], authMocks[3])

				_, err := list.Authorize(dummyAccount, dummyTransaction)
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, expectedErr)
			})
		})

		Convey("When there are commit functions", func() {
			callCounts := make([]int, len(authMocks))
			for i, authzer := range authMocks {
				commitFunc := rule.CommitFunc(nil)
				if i%2 == 0 {
					i := i // necessary for closure of separate variable
					commitFunc = func() { callCounts[i]++ }
				}
				authzer.EXPECT().
					Authorize(gomock.Eq(dummyAccount), gomock.Eq(dummyTransaction)).
					Return(commitFunc, nil)
			}
			commitFunc, _ := list.Authorize(dummyAccount, dummyTransaction)

			Convey("No commit function should be called immediately", func() {
				zeroedCallCounts := make([]int, len(authMocks))

				So(callCounts, ShouldResemble, zeroedCallCounts)
			})
			Convey("Returned commit function should call all of the internal ones", func() {
				expectedCallCounts := []int{1, 0, 1, 0, 1}

				commitFunc()
				So(callCounts, ShouldResemble, expectedCallCounts)
			})
		})
	})
}

func configureMocks(mocks []*mock_rule.MockAuthorizer, skipIndexes ...int) {
	for i, authzer := range mocks {
		if containsInt(skipIndexes, i) {
			continue
		}
		authzer.EXPECT().
			Authorize(gomock.Eq(dummyAccount), gomock.Eq(dummyTransaction)).
			Return(nil, nil)
	}
}

func configureMocksToErr(err error, mocks ...*mock_rule.MockAuthorizer) {
	for _, authzer := range mocks {
		authzer.EXPECT().
			Authorize(gomock.Eq(dummyAccount), gomock.Eq(dummyTransaction)).
			Return(nil, err)
	}
}

func containsInt(slc []int, value int) bool {
	for _, elm := range slc {
		if elm == value {
			return true
		}
	}
	return false
}
