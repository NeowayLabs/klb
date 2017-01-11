package azure_test

import (
	"fmt"
	"testing"
)

func genStorageAccountName() string {
	return fmt.Sprintf("klb-availset-tests-%d", rand.Intn(1000))
}

func testStorageAccountCreate(t *testing.T, f fixture.F) {
	storage := genStorageAccountName()
	fmt.Println(storage)
	/*nash.Run(
		f.Ctx,
		t,
		"./testdata/create_avail_set.sh",
		f.ResGroupName,
		availset,
		f.Location,
	)
	availSets := azure.NewAvailSet(f.Ctx, t, f.Session, f.ResGroupName)
	availSets.AssertExists(t, availset)
	*/
}

func testStorageAccountDelete(t *testing.T, f fixture.F) {

	storage := genStorageAccountName()
	fmt.Println(storage)
/*	nash.Run(
		f.Ctx,
		t,
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
		"./testdata/delete_avail_set.sh",
		f.ResGroupName,
		availset,
	)
	availSets.AssertDeleted(t, availset)
	*/
}
