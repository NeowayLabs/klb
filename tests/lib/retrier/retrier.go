//Package retrier implements a simple try again until
//the given context gets cancelled. Useful to tests that can have
//intermittent failures that do not represent a real failure,
//avoiding false positives.
package retrier

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
)

type Retrier struct {
	ctx context.Context
	t   *testing.T
	l   *log.Logger
}

type WorkFunc func() error

// New creates a new Retrier instance.
// The given context is used to model cancellation
// of any operation on this retrier
func New(
	ctx context.Context,
	t *testing.T,
	l *log.Logger,
) *Retrier {
	return &Retrier{
		ctx: ctx,
		t:   t,
		l:   l,
	}
}

// Run executes the given work function trying again if something
// goes wrong. On success it just returns, if the context gets
// cancelled it will call testing.T.Fatal with all the accumulated errors.
// The name parameter is used to aid the error messages.
func (r *Retrier) Run(
	name string,
	work WorkFunc,
) {
	r.l.Printf("retrier: starting retry loop for work %q", name)
	errs := retryUntilDone(r.ctx, r.l, name, work)
	if len(errs) > 0 {
		errmsgs := []string{
			"\n",
			fmt.Sprintf("work %q failed, errors in order:", name),
		}
		for i, err := range errs {
			errmsgs = append(
				errmsgs,
				fmt.Sprintf("error[%d]: %s", i, err),
			)
		}
		r.t.Fatal(r.t, strings.Join(errmsgs, "\n"))
	}
	r.l.Printf("retrier: success running work %q", name)
}

const backoff = 10 * time.Second

func retryUntilDone(
	ctx context.Context,
	l *log.Logger,
	name string,
	work WorkFunc,
) []error {
	var errs []error
	for {
		// buffered channel avoid goroutine leak
		result := make(chan error, 1)

		go func() {
			l.Printf("retrier: executing work: %s", name)
			result <- work()
		}()

		select {
		case res := <-result:
			{
				if res == nil {
					return nil
				}
				l.Printf("%s: got error: %s", name, res)
				errs = append(errs, res)
			}
		case <-ctx.Done():
			{
				l.Printf("retrier: %s: timeouted, returning all errors", name)
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
