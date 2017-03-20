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
func (lb *LoadBalancers) AssertExists(
	t *testing.T,
	name string,
	frontendipName string,
	privateIP string,
	poolname string,
) {
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
				return assertConfig(t, l, frontendipName, privateIP, poolname)
			}
		}
		return fmt.Errorf("unable to find %s in %s", name, res.Value)
	})
}

func assertConfig(
	t *testing.T,
	lb network.LoadBalancer,
	frontendipName string,
	privateIP string,
	poolname string,
) error {
	prop := lb.LoadBalancerPropertiesFormat
	if prop == nil {
		return fmt.Errorf("no properties found on lb: %+v", lb)
	}
	front_ips := prop.FrontendIPConfigurations
	if front_ips == nil {
		return fmt.Errorf("no frontend ip found in: %+v", prop)
	}

	for _, front_ip := range *front_ips {
		if front_ip.Name == nil {
			continue
		}
		if frontendipName != *front_ip.Name {
			continue
		}

		ip_prop := front_ip.FrontendIPConfigurationPropertiesFormat
		if ip_prop == nil {
			return fmt.Errorf("no ip config found in: %+v", front_ip)
		}
		if ip_prop.PrivateIPAddress == nil {
			return fmt.Errorf("no private ip found in: %+v", ip_prop)
		}
		if privateIP == *ip_prop.PrivateIPAddress {
			break
		}
	}

	return nil
}
