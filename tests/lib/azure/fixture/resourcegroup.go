package fixture

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/NeowayLabs/klb/tests/lib/retrier"
)

type ResourceGroup struct {
	client  resources.GroupsClient
	ctx     context.Context
	logger  *log.Logger
	retrier *retrier.Retrier
}

func NewResourceGroup(
	ctx context.Context,
	t *testing.T,
	s *Session,
	logger *log.Logger,
) *ResourceGroup {
	rg := &ResourceGroup{
		client:  resources.NewGroupsClient(s.SubscriptionID),
		ctx:     ctx,
		logger:  logger,
		retrier: retrier.New(ctx, t, logger),
	}
	rg.client.Authorizer = s.Token
	return rg
}

func (r *ResourceGroup) AssertExists(t *testing.T, name string) {
	r.retrier.Run("ResourceGroup.AssertExists", func() error {
		_, err := r.client.CheckExistence(name)
		return err
	})
}

func (r *ResourceGroup) AssertDeleted(t *testing.T, name string) {
	r.retrier.Run("ResourceGroup.AssertDeleted", func() error {
		_, err := r.client.Get(name)
		if err == nil {
			return fmt.Errorf("resource group: %q still exists", name)
		}
		return nil
	})
}

func (r *ResourceGroup) Create(t *testing.T, name string, location string) {
	r.retrier.Run("ResourceGroup.Create", func() error {
		_, err := r.client.CreateOrUpdate(name, resources.ResourceGroup{
			Location: &location,
		})
		return err
	})
}

func (r *ResourceGroup) Delete(t *testing.T, name string) {
	r.logger.Printf("ResourceGroup.Delete: %q", name)
	r.retrier.Run("ResourceGroup.Delete", func() error {
		_, err := r.client.Delete(name, nil)
		return err
	})
	r.logger.Printf("ResourceGroup.Delete finished")
}
