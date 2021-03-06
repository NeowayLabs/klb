package azure_test

import (
	"strings"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func genVnetName() string {
	return fixture.NewUniqueName("vnet")
}

type vnetDescription struct {
	name         string
	vnetAddr     string
	dnsAddresses []string
}

func createVNET(t *testing.T, f fixture.F, desc vnetDescription) {
	if len(desc.dnsAddresses) == 0 {
		desc.dnsAddresses = []string{"8.8.8.8"}
	}
	f.Shell.Run(
		"./testdata/create_vnet.sh",
		desc.name,
		f.ResGroupName,
		f.Location,
		desc.vnetAddr,
		strings.Join(desc.dnsAddresses, ","),
	)

}

func testVnetCreate(t *testing.T, f fixture.F) {
	vnet := genVnetName()
	nsg := genNsgName()
	subnet := genSubnetName()
	vnetAddress := "10.116.0.0/16"
	subnetAddress := "10.116.1.0/24"
	dnsAddresses := []string{"8.8.8.8", "4.4.4.4"}

	createVNET(t, f, vnetDescription{
		name:         vnet,
		vnetAddr:     vnetAddress,
		dnsAddresses: dnsAddresses,
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

	routeTable := genRouteTableName()

	f.Shell.Run(
		"./testdata/create_route_table.sh",
		routeTable,
		f.ResGroupName,
		f.Location,
	)

	f.Shell.Run(
		"./testdata/set_vnet_route_table.sh",
		vnet,
		subnet,
		f.ResGroupName,
		routeTable,
	)
	vnets := azure.NewVnet(f)
	vnets.AssertExists(t, vnet, vnetAddress, routeTable, dnsAddresses)
}

func TestVnet(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "Vnet_Create", timeout, location, testVnetCreate)
}
