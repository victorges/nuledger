package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"path"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	baseTestCasesDir = "./testcases"
	inputFileName    = "in.jsonl"
	outputFileName   = "out.jsonl"
)

func TestInputOutputCases(t *testing.T) {
	cases := listSubDirs(baseTestCasesDir)
	for _, caseName := range cases {
		inputFile := path.Join(baseTestCasesDir, caseName, inputFileName)
		outputFile := path.Join(baseTestCasesDir, caseName, outputFileName)

		input := readFile(inputFile)
		expected := readLines(readFile(outputFile))

		Convey("Authorizer handles "+caseName, t, func() {
			outputBuf := bytes.NewBuffer(nil)
			mainCore(input, outputBuf)

			output := readLines(outputBuf)

			So(len(output), ShouldEqual, len(expected))
			for i, line := range output {
				So(line, ShouldEqual, expected[i])
			}
		})
	}
}

func listSubDirs(path string) []string {
	files, err := ioutil.ReadDir(path)
	panicIfErr(err)

	subdirs := make([]string, len(files))
	for i, file := range files {
		if !file.IsDir() {
			panic("Unexpected non-directory in " + path + ": " + file.Name())
		}
		subdirs[i] = file.Name()
	}
	return subdirs
}

func readFile(path string) io.Reader {
	content, err := ioutil.ReadFile(path)
	panicIfErr(err)
	return bytes.NewReader(content)
}

func readLines(in io.Reader) []string {
	lines := []string{}
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	panicIfErr(scanner.Err())
	return lines
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
