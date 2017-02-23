package azure

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

type Vnet struct {
	client network.VirtualNetworksClient
	f      fixture.F
}

func NewVnet(f fixture.F) *Vnet {
	as := &Vnet{
		client: network.NewVirtualNetworksClient(f.Session.SubscriptionID),
		f:      f,
	}
	as.client.Authorizer = f.Session.Token
	return as
}

func validateVnetDnsServers(
	t *testing.T,
	expectedDnsServers []string,
	net network.VirtualNetwork,
) {
	if net.DhcpOptions == nil {
		return errors.New("The field DhcpOptions is nil!")
	}
	if net.DhcpOptions.DNSServers == nil {
		return errors.New("The field DNSServers is nil!")
	}

	dnsServers := *net.DhcpOptions.DNSServers
	err := fmt.Errorf("expected DNS servers[%s], got [%s]", expectedDnsServers, dnsServers)
	if len(dnsServer) != len(expectedDnsServers) {
		return err
	}
	for _, expectedDnsServer := range expectedDnsServers {
		found := false
		for dnsServer := range dnsServers {
			if dnsServer == expectedDnsServer {
				found = true
				break
			}
		}
		if !found {
			return err
		}
	}
}

// AssertExists checks if virtual network exists in the resource group.
// Fail tests otherwise.
func (vnet *Vnet) AssertExists(
	t *testing.T,
	name string,
	expectedAddress string,
	expectedRouteTable string,
	expectedDnsServers []string,
) {
	vnet.f.Retrier.Run(newID("Vnet", "AssertExists", name), func() error {
		net, err := vnet.client.Get(vnet.f.ResGroupName, name, "")
		if err != nil {
			return err
		}

		err := validateVnetDnsServers(t, expectedDnsServers, net)
		if err != nil {
			return err
		}

		if net.VirtualNetworkPropertiesFormat == nil {
			return errors.New("The field VirtualNetworkPropertiesFormat is nil!")
		}
		properties := *net.VirtualNetworkPropertiesFormat

		if properties.AddressSpace == nil {
			return errors.New("The field AddressSpace is nil!")
		}
		gotAddress := *net.VirtualNetworkPropertiesFormat.AddressSpace.AddressPrefixes
		if len(gotAddress) == 0 {
			return errors.New("Address is nil!")
		}
		if gotAddress[0] != expectedAddress {
			return errors.New("Address expected is " + expectedAddress + " but got " + gotAddress[0])
		}

		subnets := *properties.Subnets
		if len(subnets) == 0 {
			return errors.New("The field Subnets is nil!")
		}
		if subnets[0].SubnetPropertiesFormat == nil {
			return errors.New("The field SubnetPropertiesFormat is nil!")
		}
		if subnets[0].SubnetPropertiesFormat.RouteTable == nil {
			return errors.New("The field RouteTable is nil!")
		}
		if subnets[0].SubnetPropertiesFormat.RouteTable.ID == nil {
			return errors.New("The field ID is nil!")
		}
		gotRouteTable := *subnets[0].SubnetPropertiesFormat.RouteTable.ID
		if !strings.Contains(gotRouteTable, expectedRouteTable) {
			return errors.New("Vnet created with wrong route table. Expected: " + expectedRouteTable + "Actual: " + gotRouteTable)
		}

		return nil
	})
}
