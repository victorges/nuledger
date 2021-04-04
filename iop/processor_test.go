package iop_test

import (
	"bytes"
	"encoding/json"
	"nuledger/iop"
	mock_iop "nuledger/mocks/iop"
	"nuledger/model"
	"nuledger/model/violation"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

var startTime = time.Date(2021, time.March, 31, 14, 57, 55, 0, time.Local)

func TestIOProcessor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	Convey("Input/Output Processor", t, func() {
		handler := mock_iop.NewMockDataHandler(ctrl)

		in, out := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
		processor := iop.NewProcessor(in, out, handler)

		Convey("Does no-op on no input", func() {
			err := processor.Process()
			So(err, ShouldBeNil)
			So(out.Len(), ShouldEqual, 0)
		})

		Convey("Forwards read and received objects with no change", func() {
			input := iop.OperationInput{
				Account:     &model.Account{true, 1337},
				Transaction: &model.Transaction{"sketchy", 420, startTime},
			}
			output := iop.StateOutput{
				Account:    &model.Account{false, 7331},
				Violations: []violation.Code{"not-even-a-violation"},
			}

			handler.EXPECT().
				Handle(gomock.Eq(input)).
				Return(output, nil)

			err := json.NewEncoder(in).Encode(input)
			So(err, ShouldBeNil)
			err = processor.Process()

			So(err, ShouldBeNil)
			So(out.Len(), ShouldEqual, 95)
			// TODO: check output object
		})

	})
}
