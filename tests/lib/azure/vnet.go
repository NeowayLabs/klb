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
func (vnet *Vnet) AssertExists(t *testing.T, name string, address, routeTable string) {
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
		addressActual := *net.VirtualNetworkPropertiesFormat.AddressSpace.AddressPrefixes
		if addressActual != nil {
			if addressActual[0] != address {
				return errors.New("Address expected is " + address + " but actual is " + addressActual[0])
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
		routeTableVnet := *subnets[0].SubnetPropertiesFormat.RouteTable.ID
		if !strings.Contains(routeTableVnet, routeTable) {
			return errors.New("Vnet created with wrong route table. Expected: " + routeTable + "Actual: " + routeTableVnet)
		}

		return nil
	})
}
