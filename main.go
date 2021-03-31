package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/victorges/nudger/model"
)

func main() {
	var (
		in  = json.NewDecoder(os.Stdin)
		out = bufio.NewWriter(os.Stdout)
	)

	var op model.OperationInput
	for readOperation(in, &op) {
		// TODO: Call actual ledger to return state
		state := model.StateOutput{Account: op.Account}

		if err := writeState(out, state); err != nil {
			panic(err)
		}
	}
}

func readOperation(in *json.Decoder, op *model.OperationInput) bool {
	*op = model.OperationInput{}
	if err := in.Decode(op); err == io.EOF {
		return false
	} else if err != nil {
		panic(err)
	}
	return true
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
