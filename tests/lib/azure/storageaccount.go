package azure

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/storage"
	"github.com/NeowayLabs/klb/tests/lib/retrier"
)

type AvailSet struct {
	client   storage.AccountsClient
	ctx      context.Context
	resgroup string
}

func NewAvailSet(
	ctx context.Context,
	t *testing.T,
	s *Session,
	resgroup string,
) *AvailSet {
	as := &AvailSet{
		client:   storage.NewAccountsClient(s.SubscriptionID),
		ctx:      ctx,
		resgroup: resgroup,
	}
	as.client.Authorizer = s.token
	return as
}

// AssertExists checks if availability sets exists in the resource group.
// Fail tests otherwise.
func (availSet *AvailSet) AssertExists(t *testing.T, name string) {
	retrier.Run(availSet.ctx, t, getID("AssertExists", name), func() error {
		_, err := availSet.client.Get(availSet.resgroup, name)
		return err
	})
}

// AssertDeleted checks if resource was correctly deleted.
func (availSet *AvailSet) AssertDeleted(t *testing.T, name string) {
	retrier.Run(availSet.ctx, t, getID("AssertDeleted", name), func() error {
		_, err := availSet.client.Get(availSet.resgroup, name)
		if err == nil {
			return fmt.Errorf("resource %s should not exist", name)
		}
		return nil
	})
}

// Delete the availability set
func (availSet *AvailSet) Delete(t *testing.T, name string) {
	retrier.Run(availSet.ctx, t, getID("Delete", name), func() error {
		_, err := availSet.client.Delete(availSet.resgroup, name)
		return err
	})
}

func getID(method string, name string) string {
	return fmt.Sprintf("AvailSet.%s:%s", method, name)
}
