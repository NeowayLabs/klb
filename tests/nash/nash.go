package nash

import (
	"os"
	"testing"

	"github.com/NeowayLabs/nash"
)

func Exec(t *testing.T, testname string, script string) {
	sh, err := nash.New()

	if err != nil {
		t.Fatal(err)
	}

	sh.SetStdout(os.Stdout)
	sh.SetStderr(os.Stderr)
	nashPath := os.Getenv("HOME") + "/.nash"
	os.MkdirAll(nashPath, 0655)
	sh.SetDotDir(nashPath)

	err = sh.Exec(testname, script)
	if err != nil {
		t.Fatal(err)
	}
}
