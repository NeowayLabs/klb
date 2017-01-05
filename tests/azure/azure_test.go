package azure_test

import (
	"testing"

	"github.com/NeowayLabs/klb/tests/azure/fixture"
)

const location = "eastus"

func TestAzureResGroup(t *testing.T) {
	t.Run("Creation", testResourceGroupCreation)
	t.Run("Deletion", testResourceGroupDeletion)
}

func TestAzureAvailSet(t *testing.T) {
	fixture.Run(t, "AvailSetCreation", location, testAvailSetCreation)
	t.Run("Deletion", testAvailSetDeletion)
}
