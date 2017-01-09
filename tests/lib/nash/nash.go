package nash

import (
	"context"
	"os"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/retrier"
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

func Run(
	ctx context.Context,
	t *testing.T,
	scriptpath string,
	args ...string,
) {
	retrier.Run(ctx, t, "nash.Run:"+scriptpath, func() error {
		shell := Setup(t)
		return shell.ExecFile(scriptpath, args...)
	})
}
