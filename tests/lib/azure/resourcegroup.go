package azure

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/NeowayLabs/klb/tests/lib/retrier"
)

type ResourceGroup struct {
	client resources.GroupsClient
	ctx    context.Context
	logger *log.Logger
}

func NewResourceGroup(
	ctx context.Context,
	t *testing.T,
	s *Session,
	logger *log.Logger,
) *ResourceGroup {
	rg := &ResourceGroup{
		client: resources.NewGroupsClient(s.SubscriptionID),
		ctx:    ctx,
		logger: logger,
	}
	rg.client.Authorizer = s.token
	return rg
}

func (r *ResourceGroup) AssertExists(t *testing.T, name string) {
	retrier.Run(r.ctx, t, r.logger, "ResourceGroup.AssertExists", func() error {
		_, err := r.client.CheckExistence(name)
		return err
	})
}

func (r *ResourceGroup) AssertDeleted(t *testing.T, name string) {
	retrier.Run(r.ctx, t, r.logger, "ResourceGroup.AssertDeleted", func() error {
		_, err := r.client.Get(name)
		if err == nil {
			return fmt.Errorf("resource group: %q still exists", name)
		}
		return nil
	})
}

func (r *ResourceGroup) Create(t *testing.T, name string, location string) {
	retrier.Run(r.ctx, t, r.logger, "ResourceGroup.Create", func() error {
		_, err := r.client.CreateOrUpdate(name, resources.ResourceGroup{
			Location: &location,
		})
		return err
	})
}

func (r *ResourceGroup) Delete(t *testing.T, name string) {
	r.logger.Printf("ResourceGroup.Delete: %q", name)
	retrier.Run(r.ctx, t, r.logger, "ResourceGroup.Delete", func() error {
		_, err := r.client.Delete(name, nil)
		return err
	})
	r.logger.Printf("ResourceGroup.Delete finished")
}
