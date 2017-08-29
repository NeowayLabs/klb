package azure_test

import (
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func TestNIC(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "NICCreation", timeout, location, testNicCreate)
	fixture.Run(
		t,
		"NICLoadBalancerAddressPoolIntegration",
		timeout,
		location,
		testNicLoadBalancerAddressPoolIntegration,
	)
}

func testNicLoadBalancerAddressPoolIntegration(t *testing.T, f fixture.F) {

	vnet := genVnetName()
	nsg := genNsgName()
	nic := genNicName()
	subnet := genSubnetName()
	vnetCIDR := "10.66.0.0/16"
	privateIP := "10.66.1.100"
	subnetCIDR := "10.66.1.0/24"

	createVNET(t, f, vnetDescription{name: vnet, vnetAddr: vnetCIDR})
	createNSG(t, f, nsg)
	createSubnet(t, f, vnet, subnet, subnetCIDR, nsg)
	createNIC(t, f, vnet, subnet, nic, privateIP)

	nics := azure.NewNic(f)
	nics.AssertExists(t, nic, privateIP)

	const lbname = "niclb"
	const frontendIPName = "nicFrontIP"
	const lbPrivateIP = "10.66.1.150"
	const poolname = "niclbpool"

	createLoadBalancer(t, f, vnet, subnet, lbname, frontendIPName, lbPrivateIP, poolname)
	loadbalancer := azure.NewLoadBalancers(f)
	loadbalancer.AssertExists(t, lbname, frontendIPName, lbPrivateIP, poolname)

	poolID := getAddressPoolID(t, f, lbname, poolname)
	if poolID == "" {
		t.Fatal("unexpected empty pool ID")
	}
}

func getAddressPoolID(t *testing.T, f fixture.F, lbname string, poolname string) string {
	return execWithIPC(t, f, func(outfile string) {
		f.Shell.Run(
			"./testdata/alb_get_pool_id.sh",
			poolname,
			f.ResGroupName,
			lbname,
			outfile,
		)
	})
}

func testNicCreate(t *testing.T, f fixture.F) {
	vnet := genVnetName()
	nsg := genNsgName()
	nic := genNicName()
	subnet := genSubnetName()
	vnetCIDR := "10.116.0.0/16"
	privateIP := "10.116.1.100"
	subnetCIDR := "10.116.1.0/24"

	createVNET(t, f, vnetDescription{name: vnet, vnetAddr: vnetCIDR})
	createNSG(t, f, nsg)
	createSubnet(t, f, vnet, subnet, subnetCIDR, nsg)
	createNIC(t, f, vnet, subnet, nic, privateIP)

	nics := azure.NewNic(f)
	nics.AssertExists(t, nic, privateIP)
}

func createNIC(
	t *testing.T,
	f fixture.F,
	vnet string,
	subnet string,
	nic string,
	privateIP string,
) {
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

func createSubnet(
	t *testing.T,
	f fixture.F,
	vnet string,
	subnet string,
	subnetCIDR string,
	nsg string,
) {
	f.Shell.Run(
		"./testdata/create_subnet.sh",
		subnet,
		f.ResGroupName,
		vnet,
		subnetCIDR,
		nsg,
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
