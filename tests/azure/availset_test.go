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

func testAvailSetCreate(t *testing.T, f fixture.F) {
	availset := genAvailSetName()
	nash.Run(
		f.Ctx,
		t,
		f.Logger,
		"./testdata/create_avail_set.sh",
		f.ResGroupName,
		availset,
		f.Location,
	)
	availSets := azure.NewAvailSet(f.Ctx, t, f.Session, f.ResGroupName)
	availSets.AssertExists(t, availset)
}

func testAvailSetDelete(t *testing.T, f fixture.F) {
	availset := genAvailSetName()
	nash.Run(
		f.Ctx,
		t,
		f.Logger,
		"./testdata/create_avail_set.sh",
		f.ResGroupName,
		availset,
		f.Location,
	)

	availSets := azure.NewAvailSet(f.Ctx, t, f.Session, f.ResGroupName)
	availSets.AssertExists(t, availset)

	nash.Run(
		f.Ctx,
		t,
		f.Logger,
		"./testdata/delete_avail_set.sh",
		f.ResGroupName,
		availset,
	)
	availSets.AssertDeleted(t, availset)
}
