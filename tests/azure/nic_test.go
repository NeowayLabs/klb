package azure_test

import (
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

// TODO: Test add/remove load balancer address pool

func TestNic(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "Creation", timeout, location, testNicCreate)
	fixture.Run(
		t,
		"LoadBalancerAddressPoolIntegration",
		timeout,
		location,
		testNicLoadBalancerAddressPoolIntegration,
	)
}

func testNicLoadBalancerAddressPoolIntegration(t *testing.T, f fixture.F) {
	t.Skip("TODO")

	vnet := genVnetName()
	nsg := genNsgName()
	nic := genNicName()
	subnet := genSubnetName()
	vnetCIDR := "10.66.0.0/16"
	privateIP := "10.66.1.100"
	subnetCIDR := "10.66.1.0/24"

	createNIC(t, f, vnet, nsg, nic, subnet, vnetCIDR, subnetCIDR, privateIP)

	nics := azure.NewNic(f)
	nics.AssertExists(t, nic, nsg, privateIP)
}

func testNicCreate(t *testing.T, f fixture.F) {
	vnet := genVnetName()
	nsg := genNsgName()
	nic := genNicName()
	subnet := genSubnetName()
	vnetCIDR := "10.116.0.0/16"
	privateIP := "10.116.1.100"
	subnetCIDR := "10.116.1.0/24"

	createNIC(t, f, vnet, nsg, nic, subnet, vnetCIDR, subnetCIDR, privateIP)

	nics := azure.NewNic(f)
	nics.AssertExists(t, nic, nsg, privateIP)
}

func createNIC(
	t *testing.T,
	f fixture.F,
	vnet string,
	nsg string,
	nic string,
	subnet string,
	vnetCIDR string,
	subnetCIDR string,
	privateIP string,
) {
	createNSG(t, f, nsg)
	createVNet(t, f, vnetDescription{name: vnet, vnetAddr: vnetCIDR})

	f.Shell.Run(
		"./testdata/create_subnet.sh",
		subnet,
		f.ResGroupName,
		vnet,
		subnetCIDR,
		nsg,
	)

	f.Shell.Run(
		"./testdata/create_nic.sh",
		f.ResGroupName,
		nic,
		f.Location,
		vnet,
		subnet,
		privateIP,
	)
}

func createNSG(t *testing.T, f fixture.F, name string) {
	f.Shell.Run(
		"./testdata/create_nsg.sh",
		name,
		f.ResGroupName,
		f.Location,
	)
}

func genNicName() string {
	return fixture.NewUniqueName("nic")
}
