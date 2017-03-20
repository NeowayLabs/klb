package azure

import (
	"errors"
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

// AssertExists checks if load balancer exists in the resource group.
// Fail tests otherwise.
func (lb *LoadBalancers) AssertExists(t *testing.T, name string) {
	lb.f.Retrier.Run(newID("LoadBalancers", "AssertExists", name), func() error {
		res, err := lb.client.List(lb.f.ResGroupName)
		if err != nil {
			return err
		}
		if res.Value == nil {
			return errors.New("no load balancers found")
		}
		lbs := *res.Value
		for _, l := range lbs {
			if l.Name == nil {
				continue
			}
			if name == *l.Name {
				return nil
			}
		}
		return fmt.Errorf("unable to find %s in %s", name, res.Value)
	})
}
