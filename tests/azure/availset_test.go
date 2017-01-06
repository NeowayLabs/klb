package azure_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
	"github.com/NeowayLabs/klb/tests/lib/nash"
)

func genAvailSetName() string {
	return fmt.Sprintf("klb-availset-tests-%d", rand.Intn(1000))
}

func testAvailSetCreation(t *testing.T, f fixture.Fixture) {
	availset := genAvailSetName()
	nash.Run(
		t,
		"./testdata/create_avail_set.sh",
		f.ResGroupName,
		availset,
		f.Location,
	)
	availSets := azure.NewAvailSet(t, f.Session)
	availSets.AssertExists(t, availset, f.ResGroupName)
}

func testAvailSetDeletion(t *testing.T, f fixture.Fixture) {

	availset := genAvailSetName()
	nash.Run(
		t,
		"./testdata/create_avail_set.sh",
		f.ResGroupName,
		availset,
		f.Location,
	)

	availSets := azure.NewAvailSet(t, f.Session)
	availSets.AssertExists(t, availset, f.ResGroupName)

	nash.Run(
		t,
		"./testdata/delete_avail_set.sh",
		f.ResGroupName,
		availset,
	)
	availSets.AssertDeleted(t, availset, f.ResGroupName)
}
