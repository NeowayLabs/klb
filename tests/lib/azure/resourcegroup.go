package azure

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
)

type ResourceGroup struct {
	client resources.GroupsClient
}

func aborterr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func (r *ResourceGroup) AssertExists(t *testing.T, name string) {
	_, err := r.client.CheckExistence(name)

	if err != nil {
		t.Fatal(err)
		return
	}
}

func (r *ResourceGroup) AssertDeleted(t *testing.T, name string) {
	res, err := r.client.Get(name)

	if err == nil {
		t.Errorf("AssertDeleted: ResourceGroup '%s' should not exists", name)
		t.Log(res)
		r.Delete(t, name)
	}
}

func (r *ResourceGroup) Create(t *testing.T, name string, location string) {
	_, err := r.client.CreateOrUpdate(name, resources.ResourceGroup{
		Location: &location,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func (r *ResourceGroup) Delete(t *testing.T, name string) {
	_, err := r.client.Delete(name, nil)

	if err != nil {
		t.Fatal(err)
	}
}

func NewResourceGroup(t *testing.T, s *Session) *ResourceGroup {
	rg := &ResourceGroup{
		client: resources.NewGroupsClient(s.SubscriptionID),
	}
	rg.client.Authorizer = s.token
	return rg
}
