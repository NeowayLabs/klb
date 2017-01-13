package azure

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

type AvailSet struct {
	client compute.AvailabilitySetsClient
	f      fixture.F
}

func NewAvailSet(f fixture.F) *AvailSet {
	as := &AvailSet{
		client: compute.NewAvailabilitySetsClient(f.Session.SubscriptionID),
		f:      f,
	}
	as.client.Authorizer = f.Session.Token
	return as
}

// AssertExists checks if availability sets exists in the resource group.
// Fail tests otherwise.
func (av *AvailSet) AssertExists(t *testing.T, name string) {
	av.f.Retrier.Run(newID("AvailSet", "AssertExists", name), func() error {
		_, err := av.client.Get(av.f.ResGroupName, name)
		return err
	})
}

// AssertDeleted checks if resource was correctly deleted.
func (av *AvailSet) AssertDeleted(t *testing.T, name string) {
	av.f.Retrier.Run(newID("AvailSet", "AssertDeleted", name), func() error {
		_, err := av.client.Get(av.f.ResGroupName, name)
		if err == nil {
			return fmt.Errorf("resource %s should not exist", name)
		}
		return nil
	})
}

// Delete the availability set
func (av *AvailSet) Delete(t *testing.T, name string) {
	av.f.Retrier.Run(newID("AvailSet", "Delete", name), func() error {
		_, err := av.client.Delete(av.f.ResGroupName, name)
		return err
	})
}
