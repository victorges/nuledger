// Package iop contains the lower-level Input/Output processing logic of the
// application.
package iop

import (
	"encoding/json"
	"fmt"
	"io"
)

//go:generate ../gen_mocks.sh processor.go

// DataHandler provides a Handle function for processing the input read by the
// input processor and returning the output that should be written.
type DataHandler interface {
	Handle(OperationInput) (StateOutput, error)
}

// IOProcessor is the Input/Output Processor (or iop, hence the package name) of
// the application. It handles parsing the JSON objects from an input stream,
// processing them through a provided DataHandler and serializing the output
// JSON to an output stream.
type IOProcessor struct {
	in      *json.Decoder
	out     *json.Encoder
	handler DataHandler
}

// NewProcessor creates an IOProcessor that reads from the provided io.Reader,
// transforms the data through the provided DataHandler and writes the result to
// the provided io.Writer.
func NewProcessor(in io.Reader, out io.Writer, handler DataHandler) *IOProcessor {
	return &IOProcessor{
		in:      json.NewDecoder(in),
		out:     json.NewEncoder(out),
		handler: handler,
	}
}

// Process reads from the input stream, processes it with the data handler and
// writes to the output stream until either an error occurs or it reaches the
// end of the stream (io.EOF). An EOF is considered a non-error, the happy path,
// so that's the only case where a nil error will be returned by this function.
func (p *IOProcessor) Process() error {
	for {
		var op OperationInput
		if err := p.in.Decode(&op); err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Error reading operation JSON from input: %w", err)
		}

		state, err := p.handler.Handle(op)
		if err != nil {
			return fmt.Errorf("Error handling operation: %w", err)
		}

		if err := p.out.Encode(state); err != nil {
			return fmt.Errorf("Error writing state JSON to output: %w", err)
		}
	}
	return nil
}
