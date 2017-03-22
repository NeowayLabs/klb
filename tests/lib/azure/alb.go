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

type LoadBalancerProbe struct {
	Name     string
	Protocol string
	Port     string
	Interval string
	Count    string
	Path     string
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
		loadbalancer, err := lb.getLoadBalancer(t, name)
		if err != nil {
			return err
		}
		return assertConfig(t, loadbalancer, frontendipName, privateIP, poolname)
	})
}

// AssertProbeExists checks if load balancer exists and it has the given probe.
// Fail tests otherwise.
func (lb *LoadBalancers) AssertProbeExists(t *testing.T, lbname string, p LoadBalancerProbe) {
	lb.f.Retrier.Run(newID("LoadBalancers", "AssertProbeExists", p.Name), func() error {
		_, err := lb.getLoadBalancer(t, lbname)
		if err != nil {
			return err
		}
		// TODO
		return nil
	})
}

func (lb *LoadBalancers) getLoadBalancer(t *testing.T, name string) (network.LoadBalancer, error) {
	res, err := lb.client.List(lb.f.ResGroupName)
	if err != nil {
		return network.LoadBalancer{}, err
	}
	if res.Value == nil {
		return network.LoadBalancer{}, errors.New("no load balancers found")
	}
	lbs := *res.Value
	for _, l := range lbs {
		if l.Name == nil {
			continue
		}
		if name == *l.Name {
			return l, nil
		}
	}
	return network.LoadBalancer{}, fmt.Errorf("unable to find %s in %s", name, res.Value)
}

func assertFrontendIp(
	t *testing.T,
	prop network.LoadBalancerPropertiesFormat,
	frontendipName string,
	privateIP string,
) error {
	frontIPs := prop.FrontendIPConfigurations
	if frontIPs == nil {
		return fmt.Errorf("no frontend ip found in: %s", prop)
	}

	for _, frontIP := range *frontIPs {
		if frontIP.Name == nil {
			continue
		}
		if frontendipName != *frontIP.Name {
			continue
		}

		ipProp := frontIP.FrontendIPConfigurationPropertiesFormat
		if ipProp == nil {
			return fmt.Errorf("no ip config found in: %s", frontIP)
		}
		if ipProp.PrivateIPAddress == nil {
			return fmt.Errorf("no private ip found in: %s", ipProp)
		}
		if privateIP == *ipProp.PrivateIPAddress {
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
