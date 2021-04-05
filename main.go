package main

import (
	"io"
	"os"

	"nuledger/authorizer"
	"nuledger/iop"
)

var (
	stdin  io.Reader = os.Stdin
	stdout io.Writer = os.Stdout
)

func main() {
	processor := iop.NewProcessor(stdin, stdout, authorizer.NewHandler())
	if err := processor.Process(); err != nil {
		panic(err)
	}
}
