package azure_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
	"github.com/NeowayLabs/klb/tests/lib/nash"
)

func genVnetName() string {
	return fmt.Sprintf("klb-vnet-tests-%d", rand.Intn(1000))
}

func testVnetCreate(t *testing.T, f fixture.F) {
	vnet := genVnetName()
	nash.Run(
		f.Ctx,
		t,
		"./testdata/create_vnet.sh",
		vnet,
		f.ResGroupName,
		f.Location,
		"10.116.0.0/16",
	)
	vnets := azure.NewVnet(f.Ctx, t, f.Session, f.ResGroupName)
	vnets.AssertExists(t, vnet)
}

func testVnetDelete(t *testing.T, f fixture.F) {

	vnet := genVnetName()
	nash.Run(
		f.Ctx,
		t,
		"./testdata/create_vnet.sh",
		vnet,
		f.ResGroupName,
		f.Location,
		"10.116.0.0/16",
	)

	vnets := azure.NewVnet(f.Ctx, t, f.Session, f.ResGroupName)
	vnets.AssertExists(t, vnet)

	nash.Run(
		f.Ctx,
		t,
		"./testdata/delete_vnet.sh",
		vnet,
		f.ResGroupName,
	)
	vnets.AssertDeleted(t, vnet)
}
