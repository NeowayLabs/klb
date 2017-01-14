package azure

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

type Nsg struct {
	client network.SecurityGroupsClient
	f      fixture.F
}

func NewNsg(f fixture.F) *Nsg {
	as := &Nsg{
		client: network.NewSecurityGroupsClient(f.Session.SubscriptionID),
		f:      f,
	}
	as.client.Authorizer = f.Session.Token
	return as
}

// AssertExists checks if network security groups exists in the resource group.
// Fail tests otherwise.
func (nsg *Nsg) AssertExists(t *testing.T, name string) {
	nsg.f.Retrier.Run(newID("Nsg", "AssertExists", name), func() error {
		_, err := nsg.client.Get(nsg.f.ResGroupName, name, "")
		return err
	})
}
