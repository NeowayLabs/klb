package azure_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func genStorageAccountName() string {
	return fmt.Sprintf("klbstortests%d", rand.Intn(1000))
}

func testStorageAccountCreate(t *testing.T, f fixture.F) {
	t.Skip()

	storagename := genStorageAccountName()
	sku := "LRS"
	kind := "BlobStorage"
	tier := "Hot"

	execWithIPC(t, f, func(outputfile string) {
		f.Shell.Run(
			"./testdata/create_storage_account.sh",
			f.ResGroupName,
			storagename,
			f.Location,
		)
	})

	acc := azure.NewStorageAccount(f)
	acc.AssertExists(t, storagename)
}

func TestStorageAccount(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "StorageAccountCreate", timeout, location, testStorageAccountCreate)
}
