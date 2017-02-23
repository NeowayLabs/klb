package azure_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func genVnetName() string {
	return fmt.Sprintf("klb-vnet-tests-%d", rand.Intn(1000))
}

func testVnetCreate(t *testing.T, f fixture.F) {
	vnet := genVnetName()
	nsg := genNsgName()
	subnet := genSubnetName()
	vnetAddress := "10.116.0.0/16"
	subnetAddress := "10.116.1.0/24"
	dnsAddresses := []string{"8.8.8.8", "4.4.4.4"}

	f.Shell.Run(
		"./testdata/create_vnet.sh",
		vnet,
		f.ResGroupName,
		f.Location,
		vnetAddress,
		strings.Join(dnsAddresses, ","),
	)

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
	vnets.AssertExists(t, vnet, vnetAddress, routeTable)
}

func TestVnet(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "Vnet_Create", timeout, location, testVnetCreate)
}
