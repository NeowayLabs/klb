package nash

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/retrier"
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
	s.retrier.Run("Shell.Run:"+scriptpath, func() error {
		completeargs := []string{scriptpath}
		completeargs = append(completeargs, args...)

		dir, err := ioutil.TempDir("", "klb-tests")
		defer func() {
			err := os.RemoveAll(dir)
			if err != nil {
				s.t.Errorf("error removing tmp dir: %s", err)
			}
		}()
		if err != nil {
			s.t.Fatalf("unable to create tmp dir: %s", err)
		}
		log := &logWriter{logger: s.logger}
		cmd := exec.Command(scriptpath, args...)
		env := getEnvWithoutHome()
		cmd.Env = append(env, fmt.Sprintf("HOME=%s", dir))
		cmd.Stdout = log
		cmd.Stderr = log
		return cmd.Run()
	})
}

type logWriter struct {
	logger *log.Logger
}

func (l *logWriter) Write(b []byte) (int, error) {
	l.logger.Println("nash:" + strings.TrimSuffix(string(b), "\n"))
	return len(b), nil
}

func getEnvWithoutHome() []string {
	env := []string{}
	for _, v := range os.Environ() {
		if !strings.HasPrefix(v, "HOME=") {
			env = append(env, v)
		}
	}
	return env
}
