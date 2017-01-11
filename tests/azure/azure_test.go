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

func TestVnetSet(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "Vnet_Create", timeout, location, testVnetCreate)
	// fixture.Run(t, "Vnet_Delete", timeout, location, testVnetDelete)
}
