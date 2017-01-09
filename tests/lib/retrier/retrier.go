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
	"time"
)

type WorkFunc func() error

// Run executes the given work function trying again if something
// goes wrong. On success it just returns, if the context gets
// cancelled it will call testing.T.Fatal with all the accumulated errors.
// The name parameter is used to aid the error messages.
func Run(
	ctx context.Context,
	t *testing.T,
	name string,
	work WorkFunc,
) {
	errs := retryUntilDone(ctx, work)
	if len(errs) > 0 {
		errmsgs := []string{
			"\n",
			fmt.Sprintf("Work %q failed, errors in order:", name),
		}
		for i, err := range errs {
			errmsgs = append(
				errmsgs,
				fmt.Sprintf("error[%d]: %s", i, err),
			)
		}
		t.Fatal(t, strings.Join(errmsgs, "\n"))
	}
}

const backoff = 10 * time.Second

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
		// Sometimes we geet an error like: Number of write requests
		// for subscription 'x'
		// exceeded the limit of '1200' for time interval '01:00:00'.
		// Please try again after '5' minutes
		// Before trying again lets backoff a little.
		time.Sleep(backoff)
	}
}
