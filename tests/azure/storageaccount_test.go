package azure_test

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func genStorageAccountName() string {
	return fmt.Sprintf("klbstortests%d", rand.Intn(1000))
}

func testStorageAccountCreate(t *testing.T, f fixture.F) {
	genstorage := genStorageAccountName()
	f.Shell.Run(
		"./testdata/create_storage_account.sh",
		f.ResGroupName,
		genstorage,
		f.Location,
	)
	storage := os.Getenv("STORAGE_ACCOUNT_NAME")
	storAccount := azure.NewStorageAccount(f)
	storAccount.AssertExists(t, storage)
}
