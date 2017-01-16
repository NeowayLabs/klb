package azure

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

type RouteTable struct {
	client network.RouteTablesClient
	f      fixture.F
}

func NewRouteTable(f fixture.F) *RouteTable {
	as := &RouteTable{
		client: network.NewRouteTablesClient(f.Session.SubscriptionID),
		f:      f,
	}
	as.client.Authorizer = f.Session.Token
	return as
}

// AssertExists checks if a route table exists in the resource group.
// Fail tests otherwise.
func (r *RouteTable) AssertExists(t *testing.T, name string) {
	r.f.Retrier.Run(newID("RouteTable", "AssertExists", name), func() error {
		_, err := r.client.Get(r.f.ResGroupName, name, "")
		return err
	})
}
