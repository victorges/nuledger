package main

import (
	"bufio"
	"bytes"
	"fmt"
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

	Convey("Authorizer application", t, func() {
		Convey("Panics in case of error", func() {
			input, output := bytes.NewReader([]byte(`not a json`)), bytes.NewBuffer(nil)
			So(func() { testMain(input, output) }, ShouldPanic)
		})

		for _, caseName := range cases {
			Convey(fmt.Sprintf(`Correctly handles test case "%s"`, caseName), func() {
				input, expectedBuf := getTestCase(caseName)

				outputBuf := bytes.NewBuffer(nil)
				testMain(input, outputBuf)

				output, expected := readLines(outputBuf), readLines(expectedBuf)
				So(len(output), ShouldEqual, len(expected))
				for i := range expected {
					So(output[i], ShouldEqual, expected[i])
				}
			})
		}
	})
}

func testMain(in io.Reader, out io.Writer) {
	prevIn, prevOut := stdin, stdout
	defer func() { stdin, stdout = prevIn, prevOut }()

	stdin, stdout = in, out
	main()
}

func listSubDirs(path string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	subdirs := make([]string, len(files))
	for i, file := range files {
		if !file.IsDir() {
			panic("Unexpected non-directory in " + path + ": " + file.Name())
		}
		subdirs[i] = file.Name()
	}
	return subdirs
}

func getTestCase(caseName string) (input io.Reader, expectedOutput io.Reader) {
	inputFile := path.Join(baseTestCasesDir, caseName, inputFileName)
	outputFile := path.Join(baseTestCasesDir, caseName, outputFileName)

	return readFile(inputFile), readFile(outputFile)
}

func readFile(path string) io.Reader {
	content, err := ioutil.ReadFile(path)
	So(err, ShouldBeNil)
	return bytes.NewReader(content)
}

func readLines(in io.Reader) []string {
	lines := []string{}
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	So(scanner.Err(), ShouldBeNil)
	return lines
}
