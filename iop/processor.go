package iop

import (
	"bufio"
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
	out     *bufio.Writer
	handler DataHandler
}

func NewProcessor(in io.Reader, out io.Writer, handler DataHandler) *IOProcessor {
	return &IOProcessor{
		in:      json.NewDecoder(in),
		out:     bufio.NewWriter(out),
		handler: handler,
	}
}

func (p *IOProcessor) Process() error {
	for {
		var op model.OperationInput
		if err := p.in.Decode(&op); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		state, err := p.handler.Handle(op)
		if err != nil {
			return err
		}

		if err := writeState(p.out, state); err != nil {
			return err
		}
	}
	return nil
}

func writeState(out *bufio.Writer, state model.StateOutput) error {
	encoder := json.NewEncoder(out)
	if err := encoder.Encode(state); err != nil {
		return fmt.Errorf("Error writing state JSON to output: %w", err)
	}
	if _, err := fmt.Fprintln(out); err != nil {
		return fmt.Errorf("Error adding line break to output: %w", err)
	}
	if err := out.Flush(); err != nil {
		return fmt.Errorf("Error flusing output buffer: %w", err)
	}
	return nil
}
