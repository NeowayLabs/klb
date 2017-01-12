package azure

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/NeowayLabs/klb/tests/lib/retrier"
)

type AvailSet struct {
	client   compute.AvailabilitySetsClient
	ctx      context.Context
	resgroup string
	logger   *log.Logger
}

func NewAvailSet(
	ctx context.Context,
	t *testing.T,
	s *Session,
	logger *log.Logger,
	resgroup string,
) *AvailSet {
	as := &AvailSet{
		client:   compute.NewAvailabilitySetsClient(s.SubscriptionID),
		ctx:      ctx,
		resgroup: resgroup,
		logger:   logger,
	}
	as.client.Authorizer = s.token
	return as
}

// AssertExists checks if availability sets exists in the resource group.
// Fail tests otherwise.
func (availSet *AvailSet) AssertExists(t *testing.T, name string) {
	retrier.Run(availSet.ctx, t, newID("AvailSet", "AssertExists", name), func() error {
		_, err := availSet.client.Get(availSet.resgroup, name)
		return err
	})
}

// AssertDeleted checks if resource was correctly deleted.
func (availSet *AvailSet) AssertDeleted(t *testing.T, name string) {
	retrier.Run(availSet.ctx, t, newID("AvailSet", "AssertDeleted", name), func() error {
		_, err := availSet.client.Get(availSet.resgroup, name)
		if err == nil {
			return fmt.Errorf("resource %s should not exist", name)
		}
		return nil
	})
}

// Delete the availability set
func (availSet *AvailSet) Delete(t *testing.T, name string) {
	retrier.Run(availSet.ctx, t, newID("AvailSet", "Delete", name), func() error {
		_, err := availSet.client.Delete(availSet.resgroup, name)
		return err
	})
}
