package azure

import (
	"errors"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

type Subnet struct {
	client network.SubnetsClient
	f      fixture.F
}

func NewSubnet(f fixture.F) *Subnet {
	as := &Subnet{
		client: network.NewSubnetsClient(f.Session.SubscriptionID),
		f:      f,
	}
	as.client.Authorizer = f.Session.Token
	return as
}

// AssertExists checks if subnet exists in the resource group.
// Fail tests otherwise.
func (s *Subnet) AssertExists(t *testing.T, vnetName, subnetName, address, nsg string) {
	s.f.Retrier.Run(newID("Subnet", "AssertExists", subnetName), func() error {
		subnet, err := s.client.Get(s.f.ResGroupName, vnetName, subnetName, "")
		if err != nil {
			return err
		}
		if subnet.SubnetPropertiesFormat == nil {
			return errors.New("The field SubnetPropertiesFormat is nil!")
		}

		properties := *subnet.SubnetPropertiesFormat
		if properties.AddressPrefix == nil {
			return errors.New("The field AddressPrefix is nil!")
		}
		addressSubnet := *properties.AddressPrefix
		if addressSubnet != address {
			return errors.New("Subnet created with wrong Address. Expected: " + address + "Actual: " + addressSubnet)
		}
		if properties.NetworkSecurityGroup == nil || properties.NetworkSecurityGroup.ID == nil {
			return errors.New("The field NetworkSecurityGroup or NetworkSecurityGroup.ID is nil!")
		}
		nsgSubnet := *subnet.SubnetPropertiesFormat.NetworkSecurityGroup.ID
		if !strings.Contains(nsgSubnet, nsg) {
			return errors.New("Subnet created in the wrong Network security group. Expected: " + nsg + "Actual: " + nsgSubnet)
		}
		return nil
	})
}
