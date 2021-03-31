package iop

import (
	"encoding/json"
	"fmt"
	"io"
	"nuledger/model"
)

type DataHandler interface {
	Handle(model.OperationInput) (model.StateOutput, error)
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

func NewProcessor(in io.Reader, out io.Writer, handler DataHandler) *IOProcessor {
	return &IOProcessor{
		in:      json.NewDecoder(in),
		out:     json.NewEncoder(out),
		handler: handler,
	}
}

func (p *IOProcessor) Process() error {
	for {
		var op model.OperationInput
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
