package azure

import (
	"errors"
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
func (vnet *Vnet) AssertExists(t *testing.T, name string, address string) {
	vnet.f.Retrier.Run(newID("Vnet", "AssertExists", name), func() error {
		net, err := vnet.client.Get(vnet.f.ResGroupName, name, "")
		if err != nil {
			return err
		}

		addressActual := *net.VirtualNetworkPropertiesFormat.AddressSpace.AddressPrefixes
		if addressActual != nil {
			if addressActual[0] != address {
				return errors.New("Address expected is " + address + " but actual is " + addressActual[0])
			}
		} else {
			return errors.New("Address is a nil pointer")
		}

		return nil
	})
}
