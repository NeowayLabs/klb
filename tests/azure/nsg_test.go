package azure_test

import (
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func genNsgName() string {
	return fixture.NewUniqueName("nsg")
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

func TestNsg(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "Nsg_Create", timeout, location, testNsgCreate)
}
