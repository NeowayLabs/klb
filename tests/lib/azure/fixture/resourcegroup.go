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
		_, err := r.client.CreateOrUpdate(name, resources.Group{
			Location: &location,
		})
		return err
	})
}

// Delete will delete an existing resource group.
// It will NOT wait more than timeout to the resource group to be
// completely destroyed because that takes a VERY long time
// (much more then creating stuff).
//
// Exceeding the given timeout is not an error, if it happens it will check
// if the resource group is on deprovisioning state.
func (r *ResourceGroup) Delete(t *testing.T, name string) {
	r.logger.Printf("ResourceGroup.Delete: deleting %q", name)

	r.client.Delete(name, r.ctx.Done())

	resgroup, err := r.client.Get(name)
	if err != nil {
		r.logger.Printf("ResourceGroup.Delete finished")
		return
	}
	r.logger.Printf("ResourceGroup.Delete: still exists, checking if deprovisioning")
	if resgroup.Properties == nil {
		t.Fatal("ResourceGroup.Delete: resgroup does not have properties")
	}
	if resgroup.Properties.ProvisioningState == nil {
		t.Fatal("ResourceGroup.Delete: resgroup does not have ProvisioningState")
	}
	provisioningState := *resgroup.Properties.ProvisioningState
	expectedState := "Deleting"
	if *resgroup.Properties.ProvisioningState != expectedState {
		t.Fatalf("ResourceGroup.Delete: expected state[%s] got [%s]", expectedState, provisioningState)
	}
	r.logger.Printf("ResourceGroup.Delete: resgroup is deprovisioning, should be ok")
}
