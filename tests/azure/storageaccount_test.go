package azure_test

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
	"github.com/NeowayLabs/klb/tests/lib/nash"
)

func genStorageAccountName() string {
	return fmt.Sprintf("klbstortests%d", rand.Intn(1000))
}

func testStorageAccountCreate(t *testing.T, f fixture.F) {
	genstorage := genStorageAccountName()
	nash.Run(
		f.Ctx,
		t,
		"./testdata/create_storage_account.sh",
		f.ResGroupName,
		genstorage,
		f.Location,
	)
	storage := os.Getenv("STORAGE_ACCOUNT_NAME")
	//storAccount := azure.NewStorageAccount(f.Ctx, t, f.Session, f.ResGroupName, "LSR", "Storage")
	storAccount := azure.NewStorageAccount(f.Ctx, t, f.Session, f.ResGroupName)
	storAccount.AssertExists(t, storage)
}

/*
func testStorageAccountDelete(t *testing.T, f fixture.F) {

	storage := genStorageAccountName()
	nash.Run(
		f.Ctx,
		t,
		"./testdata/create_storage_account.sh",
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
}
*/
