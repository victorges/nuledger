package iop_test

import (
	"nuledger/iop"
	mock_iop "nuledger/mocks/iop"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestIOProcessor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	Convey("Input/Output Processor", t, func() {
		handler := mock_iop.NewMockDataHandler(ctrl)

		handler.EXPECT().
			Handle(gomock.AssignableToTypeOf(&iop.OperationInput{})).
			Return(iop.StateOutput{}, nil).
			AnyTimes()

		processor := iop.NewProcessor(nil, nil, handler)
		_ = processor

		x := 0
		Convey("When the integer is incremented", func() {
			x++

			Convey("The value should be greater by one", func() {
				So(x, ShouldEqual, 1)
			})

			x++
			Convey("The value should be greater by two", func() {
				So(x, ShouldEqual, 2)
			})
		})
	})
}
