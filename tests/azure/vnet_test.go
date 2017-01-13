package azure_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func genVnetName() string {
	return fmt.Sprintf("klb-vnet-tests-%d", rand.Intn(1000))
}

func testVnetCreate(t *testing.T, f fixture.F) {
	vnet := genVnetName()
	f.Shell.Run(
		"./testdata/create_vnet.sh",
		vnet,
		f.ResGroupName,
		f.Location,
		"10.116.0.0/16",
	)
	vnets := azure.NewVnet(f.Ctx, t, f.Session, f.Logger, f.ResGroupName)
	vnets.AssertExists(t, vnet)
}

func TestVnetSet(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "Vnet_Create", timeout, location, testVnetCreate)
}
