package azure_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func genNicName() string {
	return fmt.Sprintf("klb-nic-tests-%d", rand.Intn(99999999))
}

func testNicCreate(t *testing.T, f fixture.F) {
	vnet := genVnetName()
	nsg := genNsgName()
	nic := genNicName()
	subnet := genSubnetName()
	vnetAddress := "10.116.0.0/16"
	addrnic := "10.116.1.100"
	subnetAddress := "10.116.1.0/24"

	f.Shell.Run(
		"./testdata/create_nsg.sh",
		nsg,
		f.ResGroupName,
		f.Location,
	)

	createVNet(t, f, vnetDescription{
		name:     vnet,
		vnetAddr: vnetAddress,
	})

	f.Shell.Run(
		"./testdata/create_subnet.sh",
		subnet,
		f.ResGroupName,
		vnet,
		subnetAddress,
		nsg,
	)

	f.Shell.Run(
		"./testdata/create_nic.sh",
		f.ResGroupName,
		nic,
		f.Location,
		vnet,
		subnet,
		addrnic,
	)
	nics := azure.NewNic(f)

	nics.AssertExists(t, nic, nsg, addrnic)
}

func TestNic(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "Nic_Create", timeout, location, testNicCreate)
}
