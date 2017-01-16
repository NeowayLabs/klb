package azure

import (
	"errors"
	"fmt"
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
func (s *Subnet) AssertExists(t *testing.T, vnetName, subnetName string) {
	s.f.Retrier.Run(newID("Subnet", "AssertExists", subnetName), func() error {
		_, err := s.client.Get(s.f.ResGroupName, vnetName, subnetName, "")
		return err
	})
}

// AssertCreatedWithRightParameters checks if subnet was created with the right parameters in the resource group.
// Fail tests otherwise.
func (s *Subnet) AssertCreatedWithRightParameters(t *testing.T, vnetName, subnetName string) {
	s.f.Retrier.Run(newID("Subnet", "AssertCreatedWithRightParameters", subnetName), func() error {
		subnet, err := s.client.Get(s.f.ResGroupName, vnetName, subnetName, "")
		if err != nil {
			return err
		}
		fmt.Printf("Subnet %+v\n:", subnet)
		return errors.New("ERRO")
	})
}
