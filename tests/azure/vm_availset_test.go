package azure_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func genVMAvailSetName() string {
	return fmt.Sprintf("klb-availset-tests-%d", rand.Intn(1000))
}

func testVMAvailSetCreate(t *testing.T, f fixture.F) {
	availset := genVMAvailSetName()
	updatedomain := "3"
	faultdomain := "3"
	f.Shell.Run(
		"./testdata/create_vm_avail_set.sh",
		f.ResGroupName,
		availset,
		f.Location,
		updatedomain,
		faultdomain,
	)
	availSets := azure.NewAvailSet(f)
	availSets.AssertExists(t, availset)
}

func testVMAvailSetDelete(t *testing.T, f fixture.F) {
	availset := genVMAvailSetName()
	updatedomain := "3"
	faultdomain := "3"
	f.Shell.Run(
		"./testdata/create_vm_avail_set.sh",
		f.ResGroupName,
		availset,
		f.Location,
		updatedomain,
		faultdomain,
	)

	availSets := azure.NewAvailSet(f)
	availSets.AssertExists(t, availset)

	f.Shell.Run(
		"./testdata/delete_vm_avail_set.sh",
		f.ResGroupName,
		availset,
	)
	availSets.AssertDeleted(t, availset)
}

func TestVMAvailabilitySet(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "VMAvailabilitySet_Create", timeout, location, testVMAvailSetCreate)
	fixture.Run(t, "VMAvailabilitySet_Delete", timeout, location, testVMAvailSetDelete)
}
