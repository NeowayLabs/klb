package azure_test

import (
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func genSubnetName() string {
	return fixture.NewUniqueName("subnet")
}

func testSubnetCreate(t *testing.T, f fixture.F) {
	nsg := genNsgName()
	vnet := genVnetName()
	subnet := genSubnetName()
	vnetAddress := "10.116.0.0/16"
	subnetAddress := "10.116.1.0/24"

	createVNET(t, f, vnetDescription{
		name:     vnet,
		vnetAddr: vnetAddress,
	})

	f.Shell.Run(
		"./testdata/create_nsg.sh",
		nsg,
		f.ResGroupName,
		f.Location,
	)

	f.Shell.Run(
		"./testdata/create_subnet.sh",
		subnet,
		f.ResGroupName,
		vnet,
		subnetAddress,
		nsg,
	)
	subnets := azure.NewSubnet(f)
	subnets.AssertExists(t, vnet, subnet, subnetAddress, nsg)
}
func TestSubnet(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "testSubnetCreate", timeout, location, testSubnetCreate)
}
