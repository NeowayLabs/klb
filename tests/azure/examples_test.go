// +build examples

package azure_test

import (
	"context"
	"testing"
	"time"

	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
	testlog "github.com/NeowayLabs/klb/tests/lib/log"
	"github.com/NeowayLabs/klb/tests/lib/nash"
)

// TestExamples aims to just keep our examples working,
// but it does not validates if infrastructure has really been built.
func TestExamples(t *testing.T) {
	t.Parallel()
	examples := []struct {
		name    string
		script  string
		cleanup string
	}{
		{
			name:    "postgres",
			script:  "../../examples/azure/postgres/build.sh",
			cleanup: "../../examples/azure/postgres/cleanup.sh",
		},
		{
			name:    "snapshots",
			script:  "../../examples/azure/snapshots/build.sh",
			cleanup: "../../examples/azure/snapshots/cleanup.sh",
		},
		{
			name:    "backup",
			script:  "../../examples/azure/backup/build.sh",
			cleanup: "../../examples/azure/backup/cleanup.sh",
		},
	}

	timeout := time.Hour

	for _, e := range examples {
		example := e
		t.Run(example.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			env := fixture.NewSession(t).Env()
			runExample(
				ctx,
				t,
				env,
				example.name,
				example.script,
				example.cleanup,
			)
		})
	}
}

func runExample(
	ctx context.Context,
	t *testing.T,
	env []string,
	name string,
	script string,
	cleanup string,
) {
	// WHY: buffered channel avoid goroutine leak
	// WHY: examples do not support try again, can't use default retrier shell
	logger, teardown := testlog.New(t, "TestExamples-"+name)
	defer teardown()

	shell := nash.New(ctx, t, logger, env)

	result := make(chan error, 1)

	go func() {
		logger.Println("pre cleanup to ensure clean state")
		shell.RunOnce(cleanup)

		logger.Printf("executing script: %s", script)
		err := shell.RunOnce(script)

		logger.Printf("cleanup script: %s", cleanup)
		e := shell.RunOnce(cleanup)

		if e != nil {
			t.Errorf("unexpected error on cleanup: %s", e)
			t.Error("resources may have leaked")
		}

		result <- err
	}()

	select {
	case res := <-result:
		{
			logger.Println("finished executing script: %s", script)
			if res != nil {
				t.Fatalf("unexpected error: %s", res)
			}
		}
	case <-ctx.Done():
		{
			t.Fatalf("timeouted exceeded running script: %s", name)
		}
	}
	logger.Println("finished test")
}
