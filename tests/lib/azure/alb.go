package azure

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

type LoadBalancers struct {
	client network.LoadBalancersClient
	f      fixture.F
}

func NewLoadBalancers(f fixture.F) *LoadBalancers {
	as := &LoadBalancers{
		client: network.NewLoadBalancersClient(f.Session.SubscriptionID),
		f:      f,
	}
	as.client.Authorizer = f.Session.Token
	return as
}

// AssertExists checks if availability sets exists in the resource group.
// Fail tests otherwise.
func (av *LoadBalancers) AssertExists(t *testing.T, name string) {
	av.f.Retrier.Run(newID("LoadBalancers", "AssertExists", name), func() error {
		_, err := av.client.Get(av.f.ResGroupName, name)
		return err
	})
}

// AssertDeleted checks if resource was correctly deleted.
func (av *LoadBalancers) AssertDeleted(t *testing.T, name string) {
	av.f.Retrier.Run(newID("LoadBalancers", "AssertDeleted", name), func() error {
		_, err := av.client.Get(av.f.ResGroupName, name)
		if err == nil {
			return fmt.Errorf("resource %s should not exist", name)
		}
		return nil
	})
}
