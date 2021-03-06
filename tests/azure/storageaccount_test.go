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

func TestStorage(t *testing.T) {
	timeout := 30 * time.Minute
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
	fixture.Run(
		t,
		"StorageAccountCheckResourcesExistence",
		timeout,
		location,
		testStorageAccountCheckResourcesExistence,
	)
}

func testStorageAccountCreateBLOBHot(t *testing.T, f fixture.F) {
	sku := "Standard_LRS"
	tier := "Hot"
	name := genStorageAccountName()

	createStorageAccountBLOB(f, name, sku, tier)
	checkStorageBlobAccount(t, f, name, sku, tier)
}

func testStorageAccountCreateBLOBCold(t *testing.T, f fixture.F) {
	sku := "Standard_LRS"
	tier := "Cool"
	name := genStorageAccountName()

	createStorageAccountBLOB(f, name, sku, tier)
	checkStorageBlobAccount(t, f, name, sku, tier)
}

func testStorageAccountCreateStandardLRS(t *testing.T, f fixture.F) {
	sku := "Standard_LRS"
	kind := "Storage"
	name := genStorageAccountName()

	createStorageAccount(f, name, sku)
	checkStorageAccount(t, f, name, sku, kind)
}

func testStorageAccountCreatePremiumLRS(t *testing.T, f fixture.F) {
	sku := "Premium_LRS"
	kind := "Storage"
	name := genStorageAccountName()

	createStorageAccount(f, name, sku)
	checkStorageAccount(t, f, name, sku, kind)
}

type BlobStorageFixture struct {
	sku             string
	tier            string
	account         string
	container       string
	testfile        string
	testfileContent string
}

func setupBlobStorageFixture(t *testing.T, f fixture.F, sku string, tier string) (BlobStorageFixture, func()) {
	container := fixture.NewUniqueName("container")

	expectedContent := fixture.NewUniqueName("random-content")
	localfile, cleanup := setupTestFile(t, expectedContent)

	account := genStorageAccountName()
	return BlobStorageFixture{
		sku:             sku,
		tier:            tier,
		account:         account,
		container:       container,
		testfile:        localfile,
		testfileContent: expectedContent,
	}, cleanup
}

func testStorageAccountUploadFiles(t *testing.T, f fixture.F) {

	sf, cleanup := setupBlobStorageFixture(t, f, "Standard_LRS", "Cool")
	defer cleanup()

	createStorageAccountBLOB(f, sf.account, sf.sku, sf.tier)
	checkStorageBlobAccount(t, f, sf.account, sf.sku, sf.tier)

	createStorageAccountContainer(f, sf.account, sf.container)

	expectedPath := "/account/file"
	uploadFileBLOB(t, f, sf.account, sf.container, expectedPath, sf.testfile)
	contents := downloadFileBLOB(t, f, sf.account, sf.container, expectedPath)

	assert.EqualStrings(t, sf.testfileContent, contents, "downloading BLOB")
}

func testStorageAccountCheckResourcesExistence(t *testing.T, f fixture.F) {
	sku := "Standard_LRS"
	tier := "Cool"
	accountname := genStorageAccountName()

	//WHY: Testing multiple things on one test is bad but
	//     building the context to run tests is too expensive on the cloud
	testStorageAccountDontExist(t, f, accountname)
	createStorageAccountBLOB(f, accountname, sku, tier)
	checkStorageBlobAccount(t, f, accountname, sku, tier)
	testStorageAccountExists(t, f, accountname)

	containerName := "klb-test-container-exists"
	testContainerDontExist(t, f, accountname, containerName)
	createStorageAccountContainer(f, accountname, containerName)
	testContainerExists(t, f, accountname, containerName)

	remotepath := "/test/exists"
	localfile, cleanup := setupTestFile(t, "whatever")
	defer cleanup()

	testBLOBDontExist(t, f, accountname, containerName, remotepath)
	uploadFileBLOB(t, f, accountname, containerName, remotepath, localfile)
	testBLOBExists(t, f, accountname, containerName, remotepath)
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

func testStorageAccountExists(t *testing.T, f fixture.F, account string) {
	err := storageAccountExists(f, account)
	if err != nil {
		t.Fatalf("expected account[%s] to exist, error[%s]", account, err)
	}
}

func testStorageAccountDontExist(t *testing.T, f fixture.F, account string) {
	err := storageAccountExists(f, account)
	if err == nil {
		t.Fatalf("expected account[%s] to not exist", account)
	}
}

func testContainerExists(
	t *testing.T,
	f fixture.F,
	account string,
	container string,
) {
	err := containerExists(f, account, container)
	if err != nil {
		t.Fatalf("expected container[%s] to exist, error[%s]", container, err)
	}
}

func testContainerDontExist(
	t *testing.T,
	f fixture.F,
	account string,
	container string,
) {
	err := containerExists(f, account, container)
	if err == nil {
		t.Fatalf("expected container[%s] to not exist", container)
	}
}

func containerExists(
	f fixture.F,
	account string,
	container string,
) error {
	return f.Shell.RunOnce(
		"./testdata/storage_container_exists.sh",
		f.ResGroupName,
		account,
		container,
	)
}

func testBLOBExists(
	t *testing.T,
	f fixture.F,
	account string,
	container string,
	remotepath string,
) {
	err := blobExists(f, account, container, remotepath)
	if err != nil {
		t.Fatalf("expected path[%s] to exist, error[%s]", remotepath, err)
	}
}

func testBLOBDontExist(
	t *testing.T,
	f fixture.F,
	account string,
	container string,
	remotepath string,
) {
	err := blobExists(f, account, container, remotepath)
	if err == nil {
		t.Fatalf("expected path[%s] to not exist", remotepath)
	}
}

func blobExists(
	f fixture.F,
	account string,
	container string,
	remotepath string,
) error {
	return f.Shell.RunOnce(
		"./testdata/storage_blob_exists.sh",
		f.ResGroupName,
		account,
		container,
		remotepath,
	)
}

func storageAccountExists(
	f fixture.F,
	account string,
) error {
	return f.Shell.RunOnce(
		"./testdata/storage_account_exists.sh",
		f.ResGroupName,
		account,
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
	kind string,
) {
	accounts := azure.NewStorageAccounts(f)
	account := accounts.Account(t, name)

	assert.EqualStrings(t, name, account.Name, "checking name")
	assert.EqualStrings(t, f.Location, account.Location, "checking location")
	assert.EqualStrings(t, sku, account.Sku, "checking SKU")
	assert.EqualStrings(t, kind, account.Kind, "checking kind")
}

func checkStorageBlobAccount(
	t *testing.T,
	f fixture.F,
	name string,
	sku string,
	tier string,
) {
	accounts := azure.NewStorageAccounts(f)
	account := accounts.Account(t, name)

	assert.EqualStrings(t, name, account.Name, "checking name")
	assert.EqualStrings(t, f.Location, account.Location, "checking location")
	assert.EqualStrings(t, "BlobStorage", account.Kind, "checking kind")
	assert.EqualStrings(t, sku, account.Sku, "checking SKU")
	assert.EqualStrings(t, tier, account.AccessTier, "checking tier")
}
