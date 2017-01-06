package nash

import (
	"os"
	"testing"

	"github.com/NeowayLabs/nash"
)

func Setup(t *testing.T) *nash.Shell {
	shell, err := nash.New()

	if err != nil {
		t.Fatal(err)
	}

	shell.SetStdout(os.Stdout)
	shell.SetStderr(os.Stderr)
	nashPath := os.Getenv("HOME") + "/.nash"
	os.MkdirAll(nashPath, 0655)
	shell.SetDotDir(nashPath)

	return shell
}

func Run(t *testing.T, scriptpath string, args ...string) {
	shell := Setup(t)
	// FIXME: if this gets exported on nash remove this
	err := shell.ExecuteString("setting args", `ARGS = `+args2Nash(args))
	if err != nil {
		t.Fatal(err)
	}
	if err := shell.ExecFile(scriptpath); err != nil {
		t.Fatal(err)
	}
}

// FIXME: if this gets exported on nash remove this
func args2Nash(args []string) string {
	ret := "("

	for i := 0; i < len(args); i++ {
		ret += `"` + args[i] + `"`

		if i < (len(args) - 1) {
			ret += " "
		}
	}

	return ret + ")"
}
