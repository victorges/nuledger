package iop_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"nuledger/iop"
	mock_iop "nuledger/mocks/iop"
	"nuledger/model"
	"nuledger/model/violation"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

var startTime = time.Date(2021, time.March, 31, 17, 57, 55, 0, time.UTC)

func TestIOProcessor(t *testing.T) {
	Convey("Given an Input/Output Processor", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		handler := mock_iop.NewMockDataHandler(ctrl)

		in, out := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
		processor := iop.NewProcessor(in, out, handler)

		test := func(expectedOutputs int, inputs ...iop.OperationInput) ([]iop.StateOutput, error) {
			writeInputs(in, inputs...)
			err := processor.Process()
			return readOutputs(out, expectedOutputs), err
		}

		Convey("When it gets no input it should do nothing", func() {
			output, err := test(0)
			So(err, ShouldBeNil)
			So(output, ShouldBeEmpty)
		})

		Convey("When it gets an I/O error", func() {
			pipeIn, pipeOut := io.Pipe()
			defer pipeIn.Close()
			defer pipeOut.Close()

			testIOErr := func(in io.Reader, out io.Writer, expectedErr error) {
				processor = iop.NewProcessor(in, out, handler)
				output, err := test(0, iop.OperationInput{})
				So(output, ShouldBeEmpty)
				So(errors.Is(err, expectedErr), ShouldBeTrue)
			}

			Convey("It propagates input errors", func() {
				inputErr := errors.New("The input error")
				pipeOut.CloseWithError(inputErr)

				testIOErr(pipeIn, out, inputErr)
			})
			Convey("It propagates output errors", func() {
				outputErr := errors.New("The output error")
				pipeIn.CloseWithError(outputErr)

				handler.EXPECT().
					Handle(gomock.Eq(iop.OperationInput{})).
					Return(iop.StateOutput{}, nil)

				testIOErr(in, pipeOut, outputErr)
			})
		})

		Convey("When the handler returns an error", func() {
			input := iop.OperationInput{}
			expectedOut := make([]iop.StateOutput, 1)
			expectedErr := errors.New("Custom error")

			handler.EXPECT().
				Handle(gomock.Eq(input)).
				Return(expectedOut[0], nil)
			handler.EXPECT().
				Handle(gomock.Eq(input)).
				Return(iop.StateOutput{}, expectedErr)

			Convey("It should stop processing and return error", func() {
				output, err := test(1, input, input, input)
				So(output, ShouldResemble, expectedOut)
				So(err, ShouldNotBeNil)
				So(errors.Is(err, expectedErr), ShouldBeTrue)
			})
		})

		Convey("When objects are read from input and returned by handler", func() {
			input := iop.OperationInput{
				Account:     &model.Account{"", true, 1337},
				Transaction: &model.Transaction{"", "sketchy", 420, startTime},
			}
			expected := iop.StateOutput{
				Account:    &model.Account{"", false, 7331},
				Violations: []violation.Code{"not-even-a-violation"},
			}

			handler.EXPECT().
				Handle(gomock.Eq(input)).
				Return(expected, nil)

			Convey("It should pass around the exact same objects", func() {
				output, err := test(1, input)
				So(err, ShouldBeNil)
				So(output[0], ShouldResemble, expected)
			})
		})

		Convey("When multiple objects are read from input", func() {
			input := []iop.OperationInput{
				{Account: &model.Account{AvailableLimit: 42}},
				{Transaction: &model.Transaction{Amount: 23}},
			}
			expected := []iop.StateOutput{
				{Account: &model.Account{"", true, 13}},
				{Violations: []violation.Code{"surely-another-non-violation"}},
			}

			calls := make([]*gomock.Call, len(input))
			for i := range calls {
				calls[i] = handler.EXPECT().
					Handle(gomock.Eq(input[i])).
					Return(expected[i], nil)
			}
			gomock.InOrder(calls...)

			Convey("It should successfully process the multiple entries in order", func() {
				output, err := test(2, input...)
				So(err, ShouldBeNil)
				So(output, ShouldResemble, expected)
			})
		})
	})
}

func writeInputs(dest io.Writer, values ...iop.OperationInput) {
	enc := json.NewEncoder(dest)
	for _, value := range values {
		err := enc.Encode(value)
		So(err, ShouldBeNil)
	}
}

func readOutputs(src io.Reader, count int) []iop.StateOutput {
	values := make([]iop.StateOutput, 0, count)

	dec := json.NewDecoder(src)
	for {
		var value iop.StateOutput
		err := dec.Decode(&value)
		if err == io.EOF {
			break
		}
		So(err, ShouldBeNil)

		values = append(values, value)
	}
	So(len(values), ShouldEqual, count)

	return values
}
