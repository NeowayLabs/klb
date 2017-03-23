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
	Port     int32
	Interval int32
	Count    int32
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
		loadbalancer, err := lb.getLoadBalancer(t, lbname)
		if err != nil {
			return err
		}
		prop := loadbalancer.LoadBalancerPropertiesFormat
		if prop == nil {
			return fmt.Errorf("no properties found on: %s", loadbalancer)
		}
		if prop.Probes == nil {
			return fmt.Errorf("no probes found on: %s", prop)
		}

		for _, probe := range *prop.Probes {
			if probe.Name == nil {
				continue
			}
			name := *probe.Name
			lb.f.Logger.Printf("analyzing probe: %q", name)
			if name != p.Name {
				continue
			}
			probeProperties, err := checkProbePropertiesFormat(probe.ProbePropertiesFormat)
			if err != nil {
				return err
			}
			protocol := string(probeProperties.Protocol)
			port := *probeProperties.Port
			path := *probeProperties.RequestPath
			interval := *probeProperties.IntervalInSeconds
			if string(probeProperties.Protocol) != p.Protocol {
				return fmt.Errorf(
					"expected probe protocol: %q, got: %q",
					p.Protocol,
					protocol,
				)
			}
			if port != p.Port {
				return fmt.Errorf(
					"expected probe port: %d, got: %d",
					p.Port,
					port,
				)
			}
			if interval != p.Interval {
				return fmt.Errorf(
					"expected probe interval: %d, interval: %d",
					p.Interval,
					interval,
				)
			}
			if path != p.Path {
				return fmt.Errorf(
					"expected probe port: %s, got: %s",
					p.Path,
					path,
				)
			}
			return nil
		}
		return fmt.Errorf("unable to find probe: %+v on lb: %s", p, lbname)
	})
}

func checkProbePropertiesFormat(p *network.ProbePropertiesFormat) (*network.ProbePropertiesFormat, error) {
	if p == nil {
		return nil, fmt.Errorf("probe %q got no properties", p)
	}
	// Missing a struct validator here, if we lock to Go 1.8 could
	// Initialize struct with same layout and validation tags.
	if p.Port == nil {
		return nil, fmt.Errorf("absent Port on %q ", p)
	}
	if p.IntervalInSeconds == nil {
		return nil, fmt.Errorf("absent IntervalInSeconds on %q ", p)
	}
	if p.RequestPath == nil {
		if p.Protocol == "Http" {
			return nil, fmt.Errorf("absent RequestPath on %q ", p)
		}
		p.RequestPath = new(string)
	}
	return p, nil
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
