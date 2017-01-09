//Package retrier implements a simple try again until
//the given context gets cancelled. Useful to tests that can have
//intermittent failures that do not represent a real failure,
//avoiding false positives.
package retrier

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
)

type WorkFunc func() error

// Run executes the given work function trying again if something
// goes wrong. On success it just returns, if the context gets
// cancelled it will call testing.T.Fatal with all the accumulated errors.
func Run(
	ctx context.Context,
	t *testing.T,
	work WorkFunc,
) {
	errs := retryUntilDone(ctx, work)

	if len(errs) > 0 {
		errmsgs := []string{
			"\n",
			"Test failed, errors in order:",
		}

		for i, err := range errs {
			errmsgs = append(
				errmsgs,
				fmt.Sprintf("error[%d]: %s", i, err),
			)
		}

		// FIXME: What if a deferred call this and we call Fatal again ?
		t.Fatal(t, strings.Join(errmsgs, "\n"))
	}
}

func retryUntilDone(ctx context.Context, work WorkFunc) []error {
	var errs []error
	for {
		// buffered channel avoid goroutine leak
		result := make(chan error, 1)

		go func() {
			result <- work()
		}()

		select {
		case res := <-result:
			{
				if res == nil {
					return nil
				}
				errs = append(errs, res)
			}
		case <-ctx.Done():
			{
				return append(
					errs,
					errors.New("retrier timeout"),
				)
			}
		}
	}
}
