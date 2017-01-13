package azure

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/NeowayLabs/klb/tests/lib/retrier"
)

type Vnet struct {
	client   network.VirtualNetworksClient
	ctx      context.Context
	logger   *log.Logger
	retrier  *retrier.Retrier
	resgroup string
}

func NewVnet(
	ctx context.Context,
	t *testing.T,
	s *Session,
	l *log.Logger,
	resgroup string,
) *Vnet {
	as := &Vnet{
		client:   network.NewVirtualNetworksClient(s.SubscriptionID),
		ctx:      ctx,
		resgroup: resgroup,
		logger:   l,
		retrier:  retrier.New(ctx, t, l),
	}
	as.client.Authorizer = s.token
	return as
}

// AssertExists checks if availability sets exists in the resource group.
// Fail tests otherwise.
func (vnet *Vnet) AssertExists(t *testing.T, name string) {
	vnet.retrier.Run(newID("Vnet", "AssertExists", name), func() error {
		_, err := vnet.client.Get(vnet.resgroup, name, "")
		if err != nil {
		}
		return err
	})
}

// AssertDeleted checks if resource was correctly deleted.
func (vnet *Vnet) AssertDeleted(t *testing.T, name string) {
	vnet.retrier.Run(newID("Vnet", "AssertDeleted", name), func() error {
		_, err := vnet.client.Get(vnet.resgroup, name, "")
		if err == nil {
			return fmt.Errorf("resource %s should not exist", name)
		}
		return nil
	})
}

// Delete the availability set
func (vnet *Vnet) Delete(t *testing.T, name string) {
	vnet.retrier.Run(newID("Vnet", "Delete", name), func() error {
		c := make(chan struct{})
		_, err := vnet.client.Delete(vnet.resgroup, name, c)
		return err
	})
}
