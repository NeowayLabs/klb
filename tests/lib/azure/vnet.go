package azure

import (
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

// AssertExists checks if availability sets exists in the resource group.
// Fail tests otherwise.
func (vnet *Vnet) AssertExists(t *testing.T, name string) {
	vnet.f.Retrier.Run(newID("Vnet", "AssertExists", name), func() error {
		_, err := vnet.client.Get(vnet.f.ResGroupName, name, "")
		if err != nil {
		}
		return err
	})
}
