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
	env     []string
}

// New creates a new Shell instance.
// This shell instance will automatically retry until
// success or until the given context gets cancelled.
//
// You can inject the environment that will be used to
// run scripts. You should not provide a HOME, NASHPATH variables,
// each test is executed on a isolated HOME/NASHPATH.
func New(
	ctx context.Context,
	t *testing.T,
	logger *log.Logger,
	env []string,
) *Shell {
	return &Shell{
		ctx:     ctx,
		t:       t,
		retrier: retrier.New(ctx, t, logger),
		logger:  logger,
		env:     env,
	}
}

func (s *Shell) Run(
	scriptpath string,
	args ...string,
) {
	s.retrier.Run("Shell.Run:"+scriptpath, func() error {
		completeargs := []string{scriptpath}
		completeargs = append(completeargs, args...)

		homedir, err := ioutil.TempDir("", "klb-tests")
		defer func() {
			err := os.RemoveAll(homedir)
			if err != nil {
				s.t.Errorf("error removing tmp dir: %s", err)
			}
		}()
		if err != nil {
			s.t.Fatalf("unable to create tmp dir: %s", err)
		}
		nashdir := homedir + "/.nash"
		env := append(
			s.env,
			fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
			fmt.Sprintf("HOME=%s", homedir),
			fmt.Sprintf("NASHPATH=%s", nashdir),
		)
		s.installKLB(env)
		return s.newcmd(env, scriptpath, args...).Run()
	})
}

func (s *Shell) newcmd(env []string, name string, args ...string) *exec.Cmd {
	log := &logWriter{logger: s.logger}
	cmd := exec.Command(name, args...)
	cmd.Env = env
	cmd.Stdout = log
	cmd.Stderr = log
	return cmd
}

func (s *Shell) installKLB(env []string) {
	s.logger.Println("installing klb on isolated test NASHPATH")
	cmd := s.newcmd(env, "make", "install")
	cmd.Dir = klbprojectdir()
	err := cmd.Run()
	if err != nil {
		s.t.Fatalf("unable to install klb on tmp isolated NASHPATH: %s", err)
		return
	}
	s.logger.Println("successfully installed klb")

}

func klbprojectdir() string {
	return os.Getenv("GOPATH") + "/src/github.com/NeowayLabs/klb"
}

type logWriter struct {
	logger *log.Logger
}

func (l *logWriter) Write(b []byte) (int, error) {
	l.logger.Println("nash:" + strings.TrimSuffix(string(b), "\n"))
	return len(b), nil
}
