package azure

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/storage"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

type StorageAccounts struct {
	client storage.AccountsClient
	f      fixture.F
}

type StorageAccount struct {
	ID       string
	Name     string
	Type     string
	Location string
	Sku      string
	Tier     string
	Kind     string
}

func NewStorageAccounts(f fixture.F) *StorageAccounts {
	sa := &StorageAccounts{
		client: storage.NewAccountsClient(f.Session.SubscriptionID),
		f:      f,
	}
	sa.client.Authorizer = f.Session.Token
	return sa
}

// Account gets the storage account with the given name.
// Fails test in case of any errors.
func (s *StorageAccounts) Account(t *testing.T, name string) StorageAccount {
	acc, err := s.client.GetProperties(s.f.ResGroupName, name)
	if err != nil {
		t.Fatalf("Account:error[%s]", err)
	}
	if acc.ID == nil {
		t.Fatalf("Account:account[%+v] ID is nil", acc)
	}
	if acc.Name == nil {
		t.Fatalf("Account:account[%+v] Name is nil", acc)
	}
	if acc.Type == nil {
		t.Fatalf("Account:account[%+v] Type is nil", acc)
	}
	if acc.Location == nil {
		t.Fatalf("Account:account[%+v] Location is nil", acc)
	}
	if acc.Sku == nil {
		t.Fatalf("Account:account[%+v] Sku is nil", acc)
	}
	return StorageAccount{
		ID:       *acc.ID,
		Name:     *acc.Name,
		Type:     *acc.Type,
		Location: *acc.Location,
		Sku:      string(acc.Sku.Name),
		Tier:     string(acc.Sku.Tier),
		Kind:     string(acc.Kind),
	}
}
