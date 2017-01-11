package azure

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/storage"
	"github.com/NeowayLabs/klb/tests/lib/retrier"
)

type StorageAccount struct {
	client   storage.AccountsClient
	ctx      context.Context
	resgroup string
}

func NewStorageAccount(
	ctx context.Context,
	t *testing.T,
	s *Session,
	resgroup string,
) *StorageAccount {
	as := &StorageAccount{
		client:   storage.NewAccountsClient(s.SubscriptionID),
		ctx:      ctx,
		resgroup: resgroup,
	}
	as.client.Authorizer = s.token
	return as
}

// AssertExists checks if availability sets exists in the resource group.
// Fail tests otherwise.
func (storageAccount *StorageAccount) AssertExists(t *testing.T, name string) {
	retrier.Run(storageAccount.ctx, t, getIDStorageAccount("AssertExists", name), func() error {
		_, err := storageAccount.client.GetProperties(storageAccount.resgroup, name)
		return err
	})
}

// AssertDeleted checks if resource was correctly deleted.
func (storageAccount *StorageAccount) AssertDeleted(t *testing.T, name string) {
	retrier.Run(storageAccount.ctx, t, getIDStorageAccount("AssertDeleted", name), func() error {
		_, err := storageAccount.client.GetProperties(storageAccount.resgroup, name)
		if err == nil {
			return fmt.Errorf("resource %s should not exist", name)
		}
		return nil
	})
}

// Delete the availability set
func (storageAccount *StorageAccount) Delete(t *testing.T, name string) {
	retrier.Run(storageAccount.ctx, t, getIDStorageAccount("Delete", name), func() error {
		_, err := storageAccount.client.Delete(storageAccount.resgroup, name)
		return err
	})
}

func getIDStorageAccount(method string, name string) string {
	return fmt.Sprintf("StorageAccount.%s:%s", method, name)
}
