package nash

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/retrier"
	"github.com/NeowayLabs/nash"
	"github.com/NeowayLabs/nash/sh"
)

type Shell struct {
	shell  *nash.Shell
	stdout *bytes.Buffer
	stderr *bytes.Buffer
}

func New(t *testing.T) *Shell {
	shell, err := nash.New()
	if err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	shell.SetStdout(&stdout)
	shell.SetStderr(&stderr)

	nashPath := os.Getenv("HOME") + "/.nash"
	os.MkdirAll(nashPath, 0655)
	shell.SetDotDir(nashPath)

	return &Shell{
		shell:  shell,
		stdout: &stdout,
		stderr: &stderr,
	}
}

func (s *Shell) Setvar(name string, value string) {
	s.shell.Setvar(name, sh.NewStrObj(value))
}

func Run(
	ctx context.Context,
	t *testing.T,
	scriptpath string,
	args ...string,
) {
	s := New(t)
	retrier.Run(ctx, t, "nash.Run:"+scriptpath, func() error {
		err := s.shell.ExecFile(scriptpath, args...)
		if err != nil {
			return fmt.Errorf(
				"error: %s\n\nstdout:%s\n\nstderr:%s\n\n",
				err,
				s.stdout,
				s.stderr,
			)
		}
		return nil
	})
}
