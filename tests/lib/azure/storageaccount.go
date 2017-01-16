package azure

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/storage"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

type StorageAccount struct {
	client storage.AccountsClient
	f      fixture.F
}

func NewStorageAccount(f fixture.F) *StorageAccount {
	as := &StorageAccount{
		client: storage.NewAccountsClient(f.Session.SubscriptionID),
		f:      f,
	}
	as.client.Authorizer = f.Session.Token
	return as
}

// AssertExists checks if storage account exists in the resource group.
// Fail tests otherwise.
func (s *StorageAccount) AssertExists(t *testing.T, name string) {
	s.f.Retrier.Run(newID("StorageAccount", "AssertExists", name), func() error {
		_, err := s.client.GetProperties(s.f.ResGroupName, name)
		return err
	})
}
