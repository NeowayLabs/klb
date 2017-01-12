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
func (av *AvailSet) AssertExists(t *testing.T, name string) {
	retrier.Run(av.ctx, t, av.logger, getID("AssertExists", name), func() error {
		_, err := av.client.Get(av.resgroup, name)
		return err
	})
}

// AssertDeleted checks if resource was correctly deleted.
func (av *AvailSet) AssertDeleted(t *testing.T, name string) {
	retrier.Run(av.ctx, t, av.logger, getID("AssertDeleted", name), func() error {
		_, err := av.client.Get(av.resgroup, name)
		if err == nil {
			return fmt.Errorf("resource %s should not exist", name)
		}
		return nil
	})
}

// Delete the availability set
func (av *AvailSet) Delete(t *testing.T, name string) {
	retrier.Run(av.ctx, t, av.logger, getID("Delete", name), func() error {
		_, err := av.client.Delete(av.resgroup, name)
		return err
	})
}

func getID(method string, name string) string {
	return fmt.Sprintf("AvailSet.%s:%s", method, name)
}
