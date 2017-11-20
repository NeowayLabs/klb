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

func testStorageAccountCreate(t *testing.T, f fixture.F, sku string, tier string) {
	name := genStorageAccountName()

	createStorageAccount(f, name, sku)

	accounts := azure.NewStorageAccounts(f)
	account := accounts.Account(t, name)

	assert.EqualStrings(t, name, account.Name, "checking name")
	assert.EqualStrings(t, f.Location, account.Location, "checking location")
	assert.EqualStrings(t, sku, account.Sku, "checking SKU")
	assert.EqualStrings(t, "Storage", account.Kind, "checking kind")
	assert.EqualStrings(t, tier, account.Tier, "checking tier")
}

func testStorageAccountCreateStandardLRS(t *testing.T, f fixture.F) {
	testStorageAccountCreate(t, f, "Standard_LRS", "Standard")
}

func testStorageAccountCreatePremiumLRS(t *testing.T, f fixture.F) {
	testStorageAccountCreate(t, f, "Premium_LRS", "Premium")
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
}
