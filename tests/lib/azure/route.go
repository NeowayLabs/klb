package azure

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

type Route struct {
	client network.RoutesClient
	f      fixture.F
}

func NewRoute(f fixture.F) *Route {
	as := &Route{
		client: network.NewRoutesClient(f.Session.SubscriptionID),
		f:      f,
	}
	as.client.Authorizer = f.Session.Token
	return as
}

// AssertExists checks if a route exists in the resource group.
// Fail tests otherwise.
func (r *Route) AssertExists(t *testing.T, routeTable, route string) {
	r.f.Retrier.Run(newID("Route", "AssertExists", route), func() error {
		_, err := r.client.Get(r.f.ResGroupName, routeTable, route)
		return err
	})
}
