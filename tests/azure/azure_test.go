package azure_test

import "testing"

func TestAzureResGroup(t *testing.T) {
	t.Run("Creation", testResourceGroupCreation)
	t.Run("Deletion", testResourceGroupDeletion)
}

func TestAzureAvailSet(t *testing.T) {
	t.Run("Creation", testAvailSetCreation)
	t.Run("Deletion", testAvailSetDeletion)
}
