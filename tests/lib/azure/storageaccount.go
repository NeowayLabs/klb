package azure

import (
	"context"
	"log"
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/storage"
	"github.com/NeowayLabs/klb/tests/lib/retrier"
)

type StorageAccount struct {
	client   storage.AccountsClient
	ctx      context.Context
	logger   *log.Logger
	retrier  *retrier.Retrier
	resgroup string
}

func NewStorageAccount(
	ctx context.Context,
	t *testing.T,
	s *Session,
	logger *log.Logger,
	resgroup string,
) *StorageAccount {
	as := &StorageAccount{
		client:   storage.NewAccountsClient(s.SubscriptionID),
		ctx:      ctx,
		resgroup: resgroup,
		logger:   logger,
		retrier:  retrier.New(ctx, t, logger),
	}
	as.client.Authorizer = s.token
	return as
}

// AssertExists checks if availability sets exists in the resource group.
// Fail tests otherwise.
func (s *StorageAccount) AssertExists(t *testing.T, name string) {
	s.retrier.Run(newID("StorageAccount", "AssertExists", name), func() error {
		_, err := s.client.GetProperties(s.resgroup, name)
		return err
	})
}
