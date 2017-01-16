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
	address := "10.116.0.0/16"

	f.Shell.Run(
		"./testdata/create_vnet.sh",
		vnet,
		f.ResGroupName,
		f.Location,
		address,
	)
	vnets := azure.NewVnet(f)
	vnets.AssertExists(t, vnet, address)

	//azure_vnet_set_route_table($name, $resgroup, $vnet, $subnet, $routetable)

	vnetRouteTable := vnet + "routetable-test"
	f.Shell.Run(
		"./testdata/set_vnet_route_table.sh",
		vnetRouteTable,
		f.ResGroupName,
		f.Location,
		vnet,
		//subnet,
		//route_table,
	)
}

func TestVnet(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "Vnet_Create", timeout, location, testVnetCreate)
}
