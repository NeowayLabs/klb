package log

import (
	"flag"
	"log"
	"os"
	"strings"
	"testing"
)

const logsdir = "./testdata/logs"

type TearDownFunc func()

//New creates a log.Logger for the given testame.
//This will save the logs on our common logs dir
//or stdout, according to what is configured by the argument
//-logger passed to go test using -args.
//
//Example stdout: go test ./... -args -logger stdout
//Example file: go test ./... -args -logger file
//
//This is not very usual on Go tests but we have pretty
//long running tests that may take some time to run and
//being able to tail some logs saved on disk is useful
//on development.
//
//You should not use the logger instance after you call TearDownFunc.
func New(t *testing.T, testname string) (*log.Logger, TearDownFunc) {
	builder, ok := logbuilders[logger]
	if !ok {
		t.Fatalf("unknow logger: %s", logger)
	}
	return builder(t, testname)
}

func newFile(t *testing.T, testname string) (*log.Logger, TearDownFunc) {
	err := os.MkdirAll(logsdir, 0755)
	if err != nil {
		t.Fatalf("creating test logs dir: %s:", err)
	}
	logspath := strings.Join([]string{logsdir, testname + ".logs"}, "/")
	file, err := os.Create(logspath)
	if err != nil {
		t.Fatalf("error opening log file: %s", err)
	}
	logger := log.New(file, "", log.Ltime)
	return logger, func() {
		file.Close()
	}
}

func newStdout(t *testing.T, testname string) (*log.Logger, TearDownFunc) {
	logger := log.New(os.Stdout, testname, log.Ltime)
	return logger, func() {}
}

type loggerBuilder func(t *testing.T, testname string) (*log.Logger, TearDownFunc)

var logbuilders map[string]loggerBuilder
var logger string

func init() {
	flag.StringVar(&logger, "logger", "file", "test logger, valid values: 'stdout' 'file'")
	flag.Parse()

	logbuilders = map[string]loggerBuilder{
		"file":   newFile,
		"stdout": newStdout,
	}
}
