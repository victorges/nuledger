package main

import (
	"os"

	"nuledger/authorizer"
	"nuledger/iop"
)

func main() {
	processor := iop.NewProcessor(os.Stdin, os.Stdout, authorizer.NewHandler())
	if err := processor.Process(); err != nil {
		panic(err)
	}
}
