package azure_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/NeowayLabs/klb/tests/lib/assert"
	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

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
	accountname := genStorageAccountName()
	containerName := "klb-test-container"

	createStorageAccountBLOB(f, accountname, sku, tier)
	checkStorageAccount(t, f, accountname, sku, expectedTier, kind)

	createStorageAccountContainer(f, accountname, containerName)

	expectedContents := fixture.NewUniqueName("random-contents")
	expectedPath := "/test/uploadedFile"
	localfile, cleanup := setupTestFile(t, expectedContents)
	defer cleanup()

	uploadFileBLOB(t, f, accountname, containerName, expectedPath, localfile)
	contents := downloadFileBLOB(t, f, accountname, containerName, expectedPath)

	assert.EqualStrings(t, expectedContents, contents, "downloading BLOB")
}

func uploadFileBLOB(
	t *testing.T,
	f fixture.F,
	account string,
	container string,
	remotepath string,
	localpath string,
) {
	f.Shell.Run(
		"./testdata/storage_upload_blob.sh",
		f.ResGroupName,
		account,
		container,
		remotepath,
		localpath,
	)
}

func downloadFileBLOB(
	t *testing.T,
	f fixture.F,
	account string,
	container string,
	remotepath string,
) string {
	// FIXME: execWithIPC trims the contents, not a problem for now
	return execWithIPC(t, f, func(localpath string) {
		f.Shell.Run(
			"./testdata/storage_download_blob.sh",
			f.ResGroupName,
			account,
			container,
			remotepath,
			localpath,
		)
	})
}

func setupTestFile(t *testing.T, expectedContents string) (string, func()) {
	f, err := ioutil.TempFile("", "az-storage-account-tests")
	assert.NoError(t, err, "creating tmp file")

	n, err := io.Copy(f, bytes.NewBufferString(expectedContents))
	assert.NoError(t, err, "copying to tmp file")
	assert.EqualInts(t, len(expectedContents), int(n), "copying to tmp file")

	cleanup := func() {
		assert.NoError(t, os.Remove(f.Name()), "removing tmp file")
	}

	return f.Name(), cleanup
}

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

func createStorageAccountContainer(f fixture.F, account string, container string) {
	f.Shell.Run(
		"./testdata/create_storage_container.sh",
		f.ResGroupName,
		account,
		container,
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
