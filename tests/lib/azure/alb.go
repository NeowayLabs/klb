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

func assertFrontendIp(
	t *testing.T,
	prop network.LoadBalancerPropertiesFormat,
	frontendipName string,
	privateIP string,
) error {
	front_ips := prop.FrontendIPConfigurations
	if front_ips == nil {
		return fmt.Errorf("no frontend ip found in: %s", prop)
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
			return fmt.Errorf("no ip config found in: %s", front_ip)
		}
		if ip_prop.PrivateIPAddress == nil {
			return fmt.Errorf("no private ip found in: %s", ip_prop)
		}
		if privateIP == *ip_prop.PrivateIPAddress {
			return nil
		}
	}

	return fmt.Errorf("unable to find %q private ip %q at %s", frontendipName, privateIP, prop)
}

func assertBackendPool(
	t *testing.T,
	prop network.LoadBalancerPropertiesFormat,
	poolname string,
) error {
	backendpools := prop.BackendAddressPools
	if backendpools == nil {
		return fmt.Errorf("no backend pools found in: %s", prop)
	}

	for _, backendpool := range *backendpools {
		if backendpool.Name == nil {
			continue
		}
		if poolname == *backendpool.Name {
			return nil
		}
	}
	return fmt.Errorf("unable to find backend pool %q at %s", poolname, prop)
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
		return fmt.Errorf("no properties found on lb: %s", lb)
	}

	err := assertFrontendIp(t, *prop, frontendipName, privateIP)
	if err != nil {
		return err
	}

	return assertBackendPool(t, *prop, poolname)
}
