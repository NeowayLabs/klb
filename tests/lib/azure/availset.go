package azure

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/NeowayLabs/klb/tests/lib/retrier"
)

type AvailSet struct {
	client compute.AvailabilitySetsClient
	ctx    context.Context
}

func NewAvailSet(ctx context.Context, t *testing.T, s *Session) *AvailSet {
	as := &AvailSet{
		client: compute.NewAvailabilitySetsClient(s.SubscriptionID),
		ctx:    ctx,
	}
	as.client.Authorizer = s.token
	return as
}

func getID(method string, resgroup string) string {
	return fmt.Sprintf("AvailSet.%s:%s", method, resgroup)
}

// AssertExists checks if availability sets exists in the resource group.
// Fail tests otherwise.
func (availSet *AvailSet) AssertExists(t *testing.T, name, resgroup string) {
	retrier.Run(availSet.ctx, t, getID("AssertExists", name), func() error {
		_, err := availSet.client.Get(resgroup, name)
		return err
	})
}

// AssertDeleted checks if resource was correctly deleted.
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
