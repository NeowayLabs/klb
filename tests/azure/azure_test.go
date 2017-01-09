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

func TestAzureResGroup(t *testing.T) {
	t.Parallel()
	t.Run("Create", testResourceGroupCreation)
	t.Run("Deletion", testResourceGroupDeletion)
}

func TestAzureAvailSet(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "Create", timeout, location, testAvailSetCreation)
	fixture.Run(t, "Delete", timeout, location, testAvailSetDeletion)
}
