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

func (availSet *AvailSet) Check(t *testing.T, name, resgroup string) {
	_, err := availSet.client.Get(resgroup, name)

	if err != nil {
		t.Fatal(err)
	}
}

func (availSet *AvailSet) CheckDelete(t *testing.T, name, resgroup string) {
	_, err := availSet.client.Get(resgroup, name)

	if err == nil {
		availSet.Delete(t, name, resgroup)
	}
}

func (availSet *AvailSet) Delete(t *testing.T, name, resgroup string) {
	_, err := availSet.client.Delete(resgroup, name)

	if err != nil {
		t.Fatal(err)
	}
}
