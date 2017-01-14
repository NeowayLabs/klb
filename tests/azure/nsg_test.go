package azure_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func genNsgName() string {
	return fmt.Sprintf("klb-nsg-tests-%d", rand.Intn(1000))
}

func testNsgCreate(t *testing.T, f fixture.F) {
	nsg := genNsgName()
	f.Shell.Run(
		"./testdata/create_nsg.sh",
		nsg,
		f.ResGroupName,
		f.Location,
	)
	nsgs := azure.NewNsg(f)
	nsgs.AssertExists(t, nsg)
}

func TestNsgSet(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "Nsg_Create", timeout, location, testNsgCreate)
}
