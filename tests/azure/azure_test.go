package azure_test

import (
	"testing"
	"time"

	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

const (
	location = "eastus"
	timeout  = 5 * time.Minute
)

func TestResourceGroup(t *testing.T) {
	t.Parallel()
	t.Run("Create", testResourceGroupCreate)
	t.Run("Delete", testResourceGroupDelete)
}

func TestAvailabilitySet(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "AvailabilitySet_Create", timeout, location, testAvailSetCreate)
	fixture.Run(t, "AvailabilitySet_Delete", timeout, location, testAvailSetDelete)
}

func TestStorageAccount(t *testing.T) {
	t.Parallel()
	storageAccoutTimeout := 150 * time.Second
	fixture.Run(t, "StorageAccount_Create", storageAccoutTimeout, location, testStorageAccountCreate)
}
