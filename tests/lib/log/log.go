package log

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

type Logger struct {
	file *os.File
	t    *testing.T
}

const logsdir = "./testdata/logs"

//New creates a new logger for the given testame.
//This will save the logs on the our common logs dir.
//This is not very usual on Go tests but we have pretty
//long running tests that may take some time to run and
//being able to tail some logs is useful.
func New(t *testing.T, testname string) *Logger {
	err := os.MkdirAll(logsdir, 0755)
	if err != nil {
		t.Fatalf("creating test logs dir: %s:", err)
	}
	logspath := strings.Join([]string{logsdir, testname + ".logs"}, "/")
	file, err := os.Create(logspath)
	if err != nil {
		t.Fatalf("error opening log file: %s", err)
	}
	return &Logger{
		t:    t,
		file: file,
	}
}

func (l *Logger) Log(msg string, args ...interface{}) {
	_, err := l.Write([]byte(fmt.Sprintf(msg, args...) + "\n"))
	if err != nil {
		l.t.Fatalf("error logging msg: %s", err)
	}
}

func (l *Logger) Write(b []byte) (n int, err error) {
	//TODO: Not handling when b has multiple lines on it
	timestamp := time.Now().Format("15:04:05.000")
	timestamped := append([]byte(timestamp+": "), b...)
	l.file.Write(timestamped)
	l.file.Sync()
	// fake to the caller, we wrote more stuff :-)
	// ignoring errors here is not the worst thing in life
	return len(b), nil
}

func (l *Logger) Close() {
	l.file.Close()
}
