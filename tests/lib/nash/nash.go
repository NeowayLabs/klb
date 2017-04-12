package nash

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/retrier"
	"github.com/NeowayLabs/nash"
)

type Shell struct {
	ctx     context.Context
	t       *testing.T
	retrier *retrier.Retrier
	logger  *log.Logger
}

func New(
	ctx context.Context,
	t *testing.T,
	logger *log.Logger,
) *Shell {
	return &Shell{
		ctx:     ctx,
		t:       t,
		retrier: retrier.New(ctx, t, logger),
		logger:  logger,
	}
}

func (s *Shell) Run(
	scriptpath string,
	args ...string,
) {
	nashshell := newNashShell(s.t, &logWriter{
		logger: s.logger,
	})
	s.retrier.Run("Shell.Run:"+scriptpath, func() error {
		completeargs := []string{scriptpath}
		completeargs = append(completeargs, args...)
		err := nashshell.ExecFile(scriptpath, completeargs...)
		if err != nil {
			return fmt.Errorf(
				"error: %s, executing script: %s",
				err,
				scriptpath,
			)
		}
		return nil
	})
}

func newNashShell(t *testing.T, output io.Writer) *nash.Shell {
	shell, err := nash.New()
	if err != nil {
		t.Fatal(err)
	}

	shell.SetStdout(output)
	shell.SetStderr(output)

	dir, err := ioutil.TempDir("", "nash")
	if err != nil {
		t.Fatal(err)
	}
	nashPath := dir + "/.nash"
	os.MkdirAll(nashPath, 0655)
	shell.SetDotDir(nashPath)

	return shell
}

type logWriter struct {
	logger *log.Logger
}

func (l *logWriter) Write(b []byte) (int, error) {
	l.logger.Println(strings.TrimSuffix(string(b), "\n"))
	return len(b), nil
}
