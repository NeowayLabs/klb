package azure_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func genSubnetName() string {
	return fmt.Sprintf("klb-subnet-tests-%d", rand.Intn(1000))
}

func testSubnetCreate(t *testing.T, f fixture.F) {
	nsg := fmt.Sprintf("klb-nsg-tests-%d", rand.Intn(1000))
	f.Shell.Run(
		"./testdata/create_nsg.sh",
		nsg,
		f.ResGroupName,
		f.Location,
	)
	vnet := genVnetName()
	f.Shell.Run(
		"./testdata/create_vnet.sh",
		vnet,
		f.ResGroupName,
		f.Location,
		"10.116.0.0/16",
	)
	subnet := genSubnetName()
	f.Shell.Run(
		"./testdata/create_subnet.sh",
		subnet,
		f.ResGroupName,
		vnet,
		"10.116.1.0/24",
		nsg,
	)
	subnets := azure.NewSubnet(f)
	subnets.AssertExists(t, vnet, subnet)
}

func TestSubnetSet(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "Subnet_Create", timeout, location, testSubnetCreate)
}
