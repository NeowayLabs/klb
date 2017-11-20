package azure_test

import (
	"strings"
	"testing"
	"time"

	"github.com/NeowayLabs/klb/tests/lib/assert"
	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func genStorageAccountName() string {
	return strings.Replace(fixture.NewUniqueName("st"), "-", "", -1)
}

func createStorageAccount(
	f fixture.F,
	name string,
	sku string,
) {
	f.Shell.Run(
		"./testdata/create_storage_account.sh",
		f.ResGroupName,
		name,
		f.Location,
		sku,
	)
}

func createStorageAccountBLOB(
	f fixture.F,
	name string,
	sku string,
	accessTier string,
) {
	f.Shell.Run(
		"./testdata/create_storage_account_blob.sh",
		f.ResGroupName,
		name,
		f.Location,
		sku,
		accessTier,
	)
}

func checkStorageAccount(
	t *testing.T,
	f fixture.F,
	name string,
	sku string,
	tier string,
	kind string,
) {
	accounts := azure.NewStorageAccounts(f)
	account := accounts.Account(t, name)

	assert.EqualStrings(t, name, account.Name, "checking name")
	assert.EqualStrings(t, f.Location, account.Location, "checking location")
	assert.EqualStrings(t, sku, account.Sku, "checking SKU")
	assert.EqualStrings(t, kind, account.Kind, "checking kind")
	assert.EqualStrings(t, tier, account.Tier, "checking tier")
}

func testStorageAccountCreateBLOBHot(t *testing.T, f fixture.F) {
	sku := "Standard_LRS"
	tier := "Hot"
	//WHY: Because Azure is awesome
	expectedTier := "Standard"
	kind := "BlobStorage"
	name := genStorageAccountName()

	createStorageAccountBLOB(f, name, sku, tier)
	checkStorageAccount(t, f, name, sku, expectedTier, kind)
}

func testStorageAccountCreateBLOBCold(t *testing.T, f fixture.F) {
	sku := "Standard_LRS"
	tier := "Cool"
	//WHY: Because Azure is awesome
	expectedTier := "Standard"
	kind := "BlobStorage"
	name := genStorageAccountName()

	createStorageAccountBLOB(f, name, sku, tier)
	checkStorageAccount(t, f, name, sku, expectedTier, kind)
}

func testStorageAccountCreateStandardLRS(t *testing.T, f fixture.F) {
	sku := "Standard_LRS"
	kind := "Storage"
	name := genStorageAccountName()

	createStorageAccount(f, name, sku)
	checkStorageAccount(t, f, name, sku, "Standard", kind)
}

func testStorageAccountCreatePremiumLRS(t *testing.T, f fixture.F) {
	sku := "Premium_LRS"
	kind := "Storage"
	name := genStorageAccountName()

	createStorageAccount(f, name, sku)
	checkStorageAccount(t, f, name, sku, "Premium", kind)
}

func testStorageAccountUploadFiles(t *testing.T, f fixture.F) {
	sku := "Standard_LRS"
	kind := "BlobStorage"
	tier := "Cool"
	//WHY: Because Azure is awesome
	expectedTier := "Standard"
	name := genStorageAccountName()
	//container := "klb-test-container"

	createStorageAccountBLOB(f, name, sku, tier)
	checkStorageAccount(t, f, name, sku, expectedTier, kind)
}

func TestStorageAccount(t *testing.T) {
	timeout := 5 * time.Minute
	t.Parallel()
	fixture.Run(
		t,
		"StorageAccountCreateStandardLRS",
		timeout,
		location,
		testStorageAccountCreateStandardLRS,
	)
	fixture.Run(
		t,
		"StorageAccountCreatePremiumLRS",
		timeout,
		location,
		testStorageAccountCreatePremiumLRS,
	)
	fixture.Run(
		t,
		"StorageAccountCreateBLOBHot",
		timeout,
		location,
		testStorageAccountCreateBLOBHot,
	)
	fixture.Run(
		t,
		"StorageAccountCreateBLOBCold",
		timeout,
		location,
		testStorageAccountCreateBLOBCold,
	)
	fixture.Run(
		t,
		"StorageAccountUploadFiles",
		timeout,
		location,
		testStorageAccountUploadFiles,
	)
}
