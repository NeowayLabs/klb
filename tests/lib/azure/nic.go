package azure

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

type Nic struct {
	client network.InterfacesClient
	f      fixture.F
}

func NewNic(f fixture.F) *Nic {
	as := &Nic{
		client: network.NewInterfacesClient(f.Session.SubscriptionID),
		f:      f,
	}
	as.client.Authorizer = f.Session.Token
	return as
}

// AssertExists checks if virtual network exists in the resource group.
// Fail tests otherwise.
func (nic *Nic) AssertExists(t *testing.T, name string) {
	nic.f.Retrier.Run(newID("Nic", "AssertExists", name), func() error {
		_, err := nic.client.Get(nic.f.ResGroupName, name, "")
		if err != nil {
		}
		return err
	})
}
