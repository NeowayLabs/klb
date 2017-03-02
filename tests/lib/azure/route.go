package azure

import (
	"errors"
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

// checkAddressHoptypeProperties checks if address and hoptype properties exists in route.
func (r *Route) checkAddressHoptypeProperties(route network.Route, address, hoptype string) error {
	if route.RoutePropertiesFormat == nil {
		return errors.New("The field RoutePropertiesFormat is nil!")
	}

	properties := *route.RoutePropertiesFormat
	if properties.AddressPrefix == nil {
		return errors.New("The field AddressPrefix is nil!")
	}
	addressRoute := *properties.AddressPrefix
	if addressRoute != address {
		return errors.New("Route created with wrong Address. Expected: " + address + "Actual: " + addressRoute)
	}

	hoptypeRoute := string(properties.NextHopType)
	if hoptypeRoute != hoptype {
		return errors.New("Route created with wrong HopType. Expected: " + hoptype + "Actual: " + hoptypeRoute)
	}

	return nil
}

// AssertRouteExists checks if a route exists in the resource group.
// Fail tests otherwise.
func (r *Route) AssertRouteExists(t *testing.T, routeTableName, routeName, address, hoptype string) {
	r.f.Retrier.Run(newID("Route", "AssertExists", routeName), func() error {
		route, err := r.client.Get(r.f.ResGroupName, routeTableName, routeName)
		if err != nil {
			return err
		}

		return r.checkAddressHoptypeProperties(route, address, hoptype)
	})
}

// AssertVirtualApplianceRouteExists checks if a route of VirtualAppliance hop type exists in the resource group.
// Fail tests otherwise.
func (r *Route) AssertVirtualApplianceRouteExists(t *testing.T, routeTableName, routeName, address, hoptype, hopaddress string) {
	r.f.Retrier.Run(newID("Route", "AssertExists", routeName), func() error {
		route, err := r.client.Get(r.f.ResGroupName, routeTableName, routeName)
		if err != nil {
			return err
		}

		err = r.checkAddressHoptypeProperties(route, address, hoptype)
		if err != nil {
			return err
		}

		properties := *route.RoutePropertiesFormat
		if properties.NextHopIPAddress == nil {
			return errors.New("The field NextHopIPAddress is nil!")
		}
		hopaddressRoute := *properties.NextHopIPAddress
		if hopaddressRoute != hopaddress {
			return errors.New("Route created with wrong HopAddress. Expected: " + hopaddress + "Actual: " + hopaddressRoute)
		}

		return nil
	})
}
