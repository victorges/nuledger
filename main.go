package main

import (
	"io"
	"os"

	"nuledger/authorizer"
	"nuledger/iop"
)

func main() {
	mainCore(os.Stdin, os.Stdout)
}

func mainCore(in io.Reader, out io.Writer) {
	processor := iop.NewProcessor(in, out, authorizer.NewHandler())
	if err := processor.Process(); err != nil {
		panic(err)
	}
}
