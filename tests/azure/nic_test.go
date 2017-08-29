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

func assertLBBackendAddrPoolsOnNIC(
	t *testing.T,
	ipconfig azure.NicIPConfig,
	expectedPools []string,
) {
	if len(expectedPools) != len(ipconfig.LBBackendAddrPoolsIDs) {
		t.Fatalf(
			"expectedPools[%s] != gotPools[%s]",
			expectedPools,
			ipconfig.LBBackendAddrPoolsIDs,
		)
	}

	for _, expected := range expectedPools {
		got := false
		for _, pool := range ipconfig.LBBackendAddrPoolsIDs {
			if expected == pool {
				got = true
			}
		}
		if !got {
			t.Fatalf(
				"unable to find addrpool[%s] on [%s]",
				expected,
				ipconfig.LBBackendAddrPoolsIDs,
			)
		}
	}
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

	ipconfig := getIPConfig(t, f, nic)
	assertLBBackendAddrPoolsOnNIC(t, ipconfig, []string{})

	const lbname = "niclb"
	const frontendIPName = "nicFrontIP"
	const lbPrivateIP = "10.66.1.150"
	const poolname = "niclbpool"

	createLoadBalancer(t, f, vnet, subnet, lbname, frontendIPName, lbPrivateIP, poolname)
	loadbalancer := azure.NewLoadBalancers(f)
	loadbalancer.AssertExists(t, lbname, frontendIPName, lbPrivateIP, poolname)

	poolID := getLBAddressPoolID(t, f, lbname, poolname)
	if poolID == "" {
		t.Fatal("unexpected empty pool ID")
	}

	addLBAddressPoolOnNIC(t, f, nic, ipconfig.Name, poolID)

	ipconfig = getIPConfig(t, f, nic)
	assertLBBackendAddrPoolsOnNIC(t, ipconfig, []string{poolID})
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

func addLBAddressPoolOnNIC(
	t *testing.T,
	f fixture.F,
	nicName string,
	ipconfigName string,
	addrpoolID string,
) {
	f.Shell.Run(
		"./testdata/nic_add_lb_address_pool.sh",
		nicName,
		ipconfigName,
		f.ResGroupName,
		addrpoolID,
	)
}

func getIPConfig(t *testing.T, f fixture.F, nic string) azure.NicIPConfig {
	nics := azure.NewNic(f)
	ipconfigs, err := nics.GetIPConfigs(t, nic)

	if err != nil {
		t.Fatal(err)
	}

	if len(ipconfigs) != 1 {
		t.Fatalf("expected one ipconfig, got: %s", ipconfigs)
	}

	return ipconfigs[0]
}
