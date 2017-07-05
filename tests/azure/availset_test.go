package azure_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func genAvailSetName() string {
	return fmt.Sprintf("klb-availset-tests-%d", rand.Intn(9999999))
}

func testAvailSetCreate(t *testing.T, f fixture.F) {
	availset := genAvailSetName()
	f.Shell.Run(
		"./testdata/create_avail_set.sh",
		f.ResGroupName,
		availset,
		f.Location,
	)
	availSets := azure.NewAvailSet(f)
	availSets.AssertExists(t, availset)
}

func testAvailSetDelete(t *testing.T, f fixture.F) {
	availset := genAvailSetName()
	f.Shell.Run(
		"./testdata/create_avail_set.sh",
		f.ResGroupName,
		availset,
		f.Location,
	)

	availSets := azure.NewAvailSet(f)
	availSets.AssertExists(t, availset)

	f.Shell.Run(
		"./testdata/delete_avail_set.sh",
		f.ResGroupName,
		availset,
	)
	availSets.AssertDeleted(t, availset)
}

func TestAvailabilitySet(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "AvailabilitySet_Create", timeout, location, testAvailSetCreate)
	fixture.Run(t, "AvailabilitySet_Delete", timeout, location, testAvailSetDelete)
}
