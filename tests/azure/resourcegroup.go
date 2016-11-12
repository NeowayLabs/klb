package azure

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
)

type Resources struct {
	client resources.GroupsClient
}

func aborterr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func (r *Resources) Check(t *testing.T, resourceName string) {
	_, err := r.client.CheckExistence(resourceName)

	if err != nil {
		t.Fatal(err)
		return
	}
}

func (r *Resources) CheckDelete(t *testing.T, resourceName string) {
	_, err := r.client.CheckExistence(resourceName)

	if err == nil {
		r.Delete(t, resourceName)
	}
}

func (r *Resources) Delete(t *testing.T, resourceName string) {
	_, err := r.client.Delete(resourceName, nil)

	if err != nil {
		t.Fatal(err)
	}
}

func NewResources(t *testing.T, s *Session) *Resources {
	rg := &Resources{
		client: resources.NewGroupsClient(s.SubscriptionID),
	}
	rg.client.Authorizer = s.token
	return rg
}
