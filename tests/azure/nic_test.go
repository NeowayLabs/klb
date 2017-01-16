package azure_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func genNicName() string {
	return fmt.Sprintf("klb-nic-tests-%d", rand.Intn(1000))
}

func testNicCreate(t *testing.T, f fixture.F) {
	vnet := genVnetName()
	nsg := genNsgName()
	nic := genNicName()
	subnet := genSubnetName()
	vnetAddress := "10.116.0.0/16"
	address := "10.116.1.1"

	f.Shell.Run(
		"./testdata/create_vnet.sh",
		vnet,
		f.ResGroupName,
		f.Location,
		vnetAddress,
	)
	vnets := azure.NewVnet(f)

	f.Shell.Run(
		"./testdata/create_nsg.sh",
		nsg,
		f.ResGroupName,
		f.Location,
	)
	nsgs := azure.NewNsg(f)

	f.Shell.Run(
		"./testdata/create_subnet.sh",
		subnet,
		f.ResGroupName,
		vnet,
		"10.116.1.0/24",
		nsg,
		"10.116.0.0/16",
		f.Location,
	)
	subnets := azure.NewSubnet(f)

	f.Shell.Run(
		"./testdata/create_nic.sh",
		f.ResGroupName,
		nic,
		f.Location,
		vnet,

		// resgroup = $ARGS[1]
		// name     = $ARGS[2]
		// location = $ARGS[3]
		// vnet     = $ARGS[4]
		// subnet   = $ARGS[5]
		// address  = $ARGS[6]

	)
	vnets := azure.NewVnet(f)
	vnets.AssertExists(t, vnet)
}

func TestVnetSet(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "Vnet_Create", timeout, location, testVnetCreate)
}
