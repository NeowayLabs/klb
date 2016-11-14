package azure

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
)

type AvailSet struct {
	client compute.AvailabilitySetsClient
}

func NewAvailSet(t *testing.T, s *Session) *AvailSet {
	as := &AvailSet{
		client: compute.NewAvailabilitySetsClient(s.SubscriptionID),
	}

	as.client.Authorizer = s.token

	return as
}

// AssertExists checks if availability sets exists in the resource group.
// Fail tests otherwise.
func (availSet *AvailSet) AssertExists(t *testing.T, name, resgroup string) {
	_, err := availSet.client.Get(resgroup, name)

	if err != nil {
		t.Fatal(err)
	}
}

// AssertDeleted checks if resource was correctly deleted. Delete it and
// throw an error otherwise.
func (availSet *AvailSet) AssertDeleted(t *testing.T, name, resgroup string) {
	_, err := availSet.client.Get(resgroup, name)

	if err == nil {
		// resource exists
		availSet.Delete(t, name, resgroup)
		t.Fatalf("AssertDeleted: Resource %s should not exists", name)
	}
}

// Delete the availability set
func (availSet *AvailSet) Delete(t *testing.T, name, resgroup string) {
	_, err := availSet.client.Delete(resgroup, name)

	if err != nil {
		t.Fatal(err)
	}
}
