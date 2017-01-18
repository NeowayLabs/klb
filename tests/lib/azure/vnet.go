package azure

import (
	"errors"
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

// AssertExists checks if virtual network exists in the resource group.
// Fail tests otherwise.
func (vnet *Vnet) AssertExists(t *testing.T, name, expectedAddress, expectedRouteTable string) {
	vnet.f.Retrier.Run(newID("Vnet", "AssertExists", name), func() error {
		net, err := vnet.client.Get(vnet.f.ResGroupName, name, "")
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
		if gotAddress != nil {
			if gotAddress[0] != expectedAddress {
				return errors.New("Address expected is " + expectedAddress + " but got " + gotAddress[0])
			}
		} else {
			return errors.New("Address is a nil pointer")
		}

		if properties.Subnets == nil {
			return errors.New("The field Subnets is nil!")
		}
		subnets := *properties.Subnets
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
