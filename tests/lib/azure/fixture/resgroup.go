package fixture

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
)

// Fixture provides you the basic data to write your tests, enjoy :-)
type Fixture struct {
	//ResGroupName is the resource group name where
	//all resources will be created
	ResGroupName string
	//Session used to interact with the API
	Session *azure.Session
	//Location where resources are created
	Location string
}

type Test func(*testing.T, Fixture)

// Run creates a unique resource group based on testname and calls
// the given testfunc passing as argument all the resources required
// to test integration with Azure cloud.
//
// It's main purpose is to make your life easier and guarantee that
// the resource group will be destroyed when testfunc exits.
// It is a programming error to reference the created resource group
// after returning from testfunc (just like Go http handlers).
func Run(
	t *testing.T,
	testname string,
	location string,
	testfunc Test,
) {
	session := azure.NewSession(t)
	resgroup := fmt.Sprintf("klb-test-%s-%d", testname, rand.Intn(1000))

	resources := azure.NewResources(t, session)
	defer resources.Delete(t, resgroup)

	resources.Create(t, resgroup, location)
	resources.AssertExists(t, resgroup)

	testfunc(t, Fixture{
		ResGroupName: resgroup,
		Session:      session,
		Location:     location,
	})
}
