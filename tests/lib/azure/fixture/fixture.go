package fixture

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	testlog "github.com/NeowayLabs/klb/tests/lib/log"
	"github.com/NeowayLabs/klb/tests/lib/nash"
	"github.com/NeowayLabs/klb/tests/lib/retrier"
)

// Fixture provides you the basic data to write your tests, enjoy :-)
type F struct {
	//Ctx created with the timeout provided for this test.
	//Use it to properly timeout your tests
	Ctx context.Context
	//ResGroupName is the resource group name where
	//all resources will be created
	ResGroupName string
	//Session used to interact with the Azure API
	Session *Session
	//Location where resources are created
	Location string
	//Name is the test name
	Name string
	//Logger useful to log on your tests, bypass go test default
	Logger *log.Logger
	//Shell nash shell wrapper, ready to execute scripts
	Shell *nash.Shell
	//Retrier retrier can be used to run functions until context is cancelled
	Retrier *retrier.Retrier
}

type Test func(*testing.T, F)

const resourceCleanupTimeout = 3 * time.Minute

// Run creates a unique resource group based on testname and calls
// the given testfunc passing as argument all the resources required
// to test integration with Azure cloud.
//
// It's main purpose is to make your life easier and guarantee that
// the resource group will be destroyed when testfunc exits.
//
// When a resource group is destroyed all resources inside it are
// destroyed too, so you don't need to worry with any cleanup inside your
// tests, just have fun.
//
// Since testing is pretty slow and the fixture guarantee unique named
// resource groups it will also run the your test in parallel with others
// so you don't have to die waiting for a result.
//
// It is a programming error to reference the created resource group
// after returning from testfunc (just like Go http handlers).
func Run(
	t *testing.T,
	testname string,
	timeout time.Duration,
	location string,
	testfunc Test,
) {
	//FIXME: We could remove testname on Go 1.8
	t.Run(testname, func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		logger, teardown := testlog.New(t, testname)
		defer teardown()

		session := NewSession(t)
		resgroup := fmt.Sprintf("klb-test-fixture-%s-%d", testname, rand.Intn(9999999))

		resources := NewResourceGroup(ctx, t, session, logger)
		defer func() {
			// We cant use an expired context when cleaning up state from Azure.
			ctx, cancel := context.WithTimeout(context.Background(), resourceCleanupTimeout)
			defer cancel()
			resources := NewResourceGroup(ctx, t, session, logger)
			resources.Delete(t, resgroup)
		}()

		logger.Printf("fixture: setting up resgroup %q at %q", resgroup, location)
		resources.Create(t, resgroup, location)
		resources.AssertExists(t, resgroup)
		logger.Printf("fixture: created resgroup %q with success", resgroup)

		logger.Println("fixture: calling test function")
		testfunc(t, F{
			Ctx:          ctx,
			Name:         testname,
			ResGroupName: resgroup,
			Session:      session,
			Location:     location,
			Logger:       logger,
			Shell:        nash.New(ctx, t, logger),
			Retrier:      retrier.New(ctx, t, logger),
		})
		logger.Printf("fixture: finished, failed=%t", t.Failed())
	})
}
