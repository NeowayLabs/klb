package fixture

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/NeowayLabs/klb/tests/lib/azure"
)

// Fixture provides you the basic data to write your tests, enjoy :-)
type F struct {
	//ResGroupName is the resource group name where
	//all resources will be created
	ResGroupName string
	//Session used to interact with the API
	Session *azure.Session
	//Location where resources are created
	Location string
	//Ctx created with the timeout provided for this test.
	//Use it to properly timeout your tests
	Ctx context.Context
}

type Test func(*testing.T, F)

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
	t.Run(testname, func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		session := azure.NewSession(t)
		resgroup := fmt.Sprintf("klb-test-fixture-%s-%d", testname, rand.Intn(9999999))

		resources := azure.NewResourceGroup(ctx, t, session)
		defer func() {
			// We cant use an expired context when cleaning up state from Azure.
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			resources := azure.NewResourceGroup(ctx, t, session)
			resources.Delete(t, resgroup)
		}()

		fmt.Println("NOME", resgroup)

		resources.Create(t, resgroup, location)
		resources.AssertExists(t, resgroup)

		testfunc(t, F{
			ResGroupName: resgroup,
			Session:      session,
			Location:     location,
			Ctx:          ctx,
		})
	})
}
