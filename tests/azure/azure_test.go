package azure_test

import (
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

const location = "eastus"

func TestAzureResGroup(t *testing.T) {
	t.Run("Create", testResourceGroupCreation)
	t.Run("Deletion", testResourceGroupDeletion)
}

func TestAzureAvailSet(t *testing.T) {
	fixture.Run(t, "Create", location, testAvailSetCreation)
	fixture.Run(t, "Delete", location, testAvailSetDeletion)
}
