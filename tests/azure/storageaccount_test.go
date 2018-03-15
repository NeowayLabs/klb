package azure_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/NeowayLabs/klb/tests/lib/assert"
	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func TestStorage(t *testing.T) {
	timeout := 15 * time.Minute
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
	fixture.Run(
		t,
		"BlobFSUploadsWhenAccountAndContainerExists",
		timeout,
		location,
		testBlobFSUploadsWhenAccountAndContainerExists,
	)
	testBlobFSDownloadDir(t, timeout, location)
	testBlobFSUploadDir(t, timeout, location)
	testBlobFSListFiles(t, timeout, location)
	testBlobFSListDirs(t, timeout, location)
	testBlobFSCreatesAccountAndContainerIfNonExistent(
		t,
		timeout,
		location,
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

func checkBlobFSUploadDir(t *testing.T, f fixture.F, remotedir string) {
	account := genStorageAccountName()
	container := fixture.NewUniqueName("uploadir")

	type TestFile struct {
		path    string
		content string
	}

	createFilesOnDir := func(dir string, filesCount int) []TestFile {
		assert.NoError(t, os.MkdirAll(dir, 0644))
		files := []TestFile{}
		for i := 0; i < filesCount; i++ {
			localfile := filepath.Join(dir, strconv.Itoa(i))
			content := fixture.NewUniqueName("somecontent")
			ioutil.WriteFile(localfile, []byte(content), 0644)
			files = append(files, TestFile{
				path:    localfile,
				content: content,
			})
		}
		return files
	}

	tdir, err := ioutil.TempDir("", "uploader-test")
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, os.RemoveAll(tdir))
	}()

	localfiles := []TestFile{}
	localfiles = append(localfiles, createFilesOnDir(
		filepath.Join(tdir), 2)...)
	localfiles = append(localfiles, createFilesOnDir(
		filepath.Join(tdir, "level1"), 1)...)
	localfiles = append(localfiles, createFilesOnDir(
		filepath.Join(tdir, "level1", "level2"), 3)...)

	sku := "Standard_LRS"
	tier := "Cool"

	fs := newBlobFS(account, sku, tier, container)
	fs.UploadDir(t, f, remotedir, tdir)

	checkStorageBlobAccount(t, f, account, sku, tier)

	expectedRemoteFiles := []string{}
	remoteToLocalFile := map[string]TestFile{}

	for _, file := range localfiles {
		remotepath := filepath.Join(
			remotedir,
			strings.TrimPrefix(file.path, tdir))
		remoteToLocalFile[remotepath] = file
		expectedRemoteFiles = append(expectedRemoteFiles, remotepath)
	}

	for remotepath, file := range remoteToLocalFile {
		filecontent := downloadFileBLOB(t, f, account, container, remotepath)
		assert.EqualStrings(
			t,
			file.content,
			filecontent,
			"checking uploaded BLOB individually")
	}
}

func testBlobFSDownloadDir(
	t *testing.T,
	timeout time.Duration,
	location string,
) {
	type TestCase struct {
		name        string
		remoteFiles []string
		wantedFiles []string
		downloadDir string
	}

	// TODO: add more scenarios
	tests := []TestCase{
		{
			name: "Root",
			remoteFiles: []string{
				"/test1",
				"/test2",
				"/test3",
			},
			downloadDir: "/",
			wantedFiles: []string{
				"/test1",
				"/test2",
				"/test3",
			},
		},
	}

	for _, test := range tests {
		name := "BlobFSDownloadDir" + test.name
		remoteFiles := test.remoteFiles
		wantedFiles := test.wantedFiles
		downloadDir := test.downloadDir

		fixture.Run(t, name, timeout, location, func(t *testing.T, f fixture.F) {
			fs := setupFS()

			createStorageAccountBLOB(f, fs.account, fs.sku, fs.tier)
			checkStorageBlobAccount(t, f, fs.account, fs.sku, fs.tier)
			createStorageAccountContainer(f, fs.account, fs.container)

			uploadedFiles := map[string]string{}

			for _, remoteFile := range remoteFiles {
				expectedContent := fixture.NewUniqueName("random-content")
				localFile, cleanup := setupTestFile(t, expectedContent)
				fs.Upload(t, f, remoteFile, localFile)
				cleanup()

				if filepath.Base(remoteFile) == downloadDir {
					uploadedFiles[remoteFile] = expectedContent
				}
			}

			tmpdir, err := ioutil.TempDir("", "az-download-dir")
			assert.NoError(t, err)
			defer func() {
				assert.NoError(t, os.RemoveAll(tmpdir))
			}()

			fs.DownloadDir(t, f, tmpdir, downloadDir)

			gotFiles := []string{}

			err = filepath.Walk(tmpdir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					t.Fatalf("fatal error[%s] walking download dir [%s]", err, tmpdir)
				}
				if info.IsDir() {
					return nil
				}
				gotFiles = append(gotFiles, strings.TrimPrefix(path, tmpdir))
				return nil
			})
			assert.NoError(t, err)

			if len(wantedFiles) != len(gotFiles) {
				t.Fatalf("wanted files[%s] got[%s]", wantedFiles, gotFiles)
			}

			for _, wantedFile := range wantedFiles {
				got := false
				for _, gotFile := range gotFiles {
					if gotFile == wantedFile {
						got = true
					}
				}
				if !got {
					t.Fatalf("wanted files[%s] got[%s]", wantedFiles, gotFiles)
				}
				// TODO: Validate file contents too
			}
		})
	}

}

func testBlobFSUploadDir(
	t *testing.T,
	timeout time.Duration,
	location string,
) {
	type TestCase struct {
		name      string
		remotedir string
	}

	tests := []TestCase{
		{name: "Root", remotedir: "/"},
		{name: "OneLevel", remotedir: "/remote1"},
		{name: "TwoLevels", remotedir: "/remote1/remote2"},
	}

	for _, test := range tests {
		name := "BlobFSUploadDir" + test.name
		remotedir := test.remotedir
		fixture.Run(t, name, timeout, location, func(t *testing.T, f fixture.F) {
			checkBlobFSUploadDir(t, f, remotedir)
		})
	}

}

type listTestOperation struct {
	dir           string
	expectedPaths []string
}

type listTestCase struct {
	name        string
	remotefiles []string
	ops         []listTestOperation
}

func setupFS() blobFS {
	account := genStorageAccountName()
	container := fixture.NewUniqueName("list")
	sku := "Standard_LRS"
	tier := "Cool"
	return newBlobFS(account, sku, tier, container)
}

func testBlobFSList(
	t *testing.T,
	testprefix string,
	timeout time.Duration,
	location string,
	list func(*testing.T, fixture.F, blobFS, string) []string,
	tests []listTestCase,
) {

	checkDir := func(
		t *testing.T,
		f fixture.F,
		fs blobFS,
		remotedir string,
		wantpaths []string,
	) {
		gotpaths := list(t, f, fs, remotedir)
		if len(gotpaths) != len(wantpaths) {
			t.Fatalf(
				"want files[%s] size[%d] got[%s] size[%d]",
				wantpaths,
				len(wantpaths),
				gotpaths,
				len(gotpaths),
			)
		}

		for _, wantpath := range wantpaths {
			found := false
			for _, gotfile := range gotpaths {
				if wantpath == gotfile {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("unable to find [%s] on [%s]", wantpath, gotpaths)
			}
		}
	}

	for _, _t := range tests {
		test := _t
		tname := "BlobFS" + testprefix + test.name
		fixture.Run(t, tname, timeout, location, func(t *testing.T, f fixture.F) {

			fs := setupFS()

			localfile, cleanup := setupTestFile(t, "whatever")
			defer cleanup()

			createStorageAccountBLOB(f, fs.account, fs.sku, fs.tier)
			checkStorageBlobAccount(t, f, fs.account, fs.sku, fs.tier)
			createStorageAccountContainer(f, fs.account, fs.container)

			for _, remotefile := range test.remotefiles {
				fs.Upload(t, f, remotefile, localfile)
			}

			for _, listoperation := range test.ops {
				checkDir(t, f, fs, listoperation.dir, listoperation.expectedPaths)
			}
		})
	}
}

func testBlobFSListDirs(t *testing.T, timeout time.Duration, location string) {
	testBlobFSList(
		t,
		"ListDirs",
		timeout,
		location,
		func(t *testing.T, f fixture.F, fs blobFS, remotepath string) []string {
			return fs.ListDir(t, f, remotepath)
		},
		[]listTestCase{
			{
				name:        "OneDir",
				remotefiles: []string{"/test/file"},
				ops: []listTestOperation{
					{
						dir:           "/",
						expectedPaths: []string{"/test"},
					},
				},
			},
			{
				name: "MultipleFiles",
				remotefiles: []string{
					"/test/file1",
					"/test/file2",
					"/test/file3",
				},
				ops: []listTestOperation{
					{
						dir: "/",
						expectedPaths: []string{
							"/test",
						},
					},
				},
			},
			{
				name: "MultipleDirs",
				remotefiles: []string{
					"/test1/file1",
					"/test2/file2",
					"/test3/file3",
				},
				ops: []listTestOperation{
					{
						dir: "/",
						expectedPaths: []string{
							"/test1",
							"/test2",
							"/test3",
						},
					},
				},
			},
			{
				// WHY: The alternative is to implemented a full hierarchical FS
				// on top of azure blob, seems excessive.
				name: "PathCanBeFileAndDir",
				remotefiles: []string{
					"/test/dir",
					"/test/dir/file",
				},
				ops: []listTestOperation{
					{
						dir:           "/test",
						expectedPaths: []string{"/test/dir"},
					},
				},
			},
			{
				name:        "EmptyDir",
				remotefiles: []string{"/test"},
				ops: []listTestOperation{
					{
						dir:           "/",
						expectedPaths: []string{},
					},
				},
			},
			{
				name:        "WrongDir",
				remotefiles: []string{},
				ops: []listTestOperation{
					{
						dir:           "/",
						expectedPaths: []string{},
					},
					{
						dir:           "/test",
						expectedPaths: []string{},
					},
				},
			},
			{
				name: "NestedDirs",
				remotefiles: []string{
					"/test/file",
					"/test/dir1/file",
					"/test/dir2/file",
					"/test/dir2/dir3/file1",
					"/test/dir2/dir3/file2",
					"/test/dir2/dir3/dir4/file1",
				},
				ops: []listTestOperation{
					{
						dir: "/test",
						expectedPaths: []string{
							"/test/dir1",
							"/test/dir2",
						},
					},
					{
						dir:           "/test/dir2",
						expectedPaths: []string{"/test/dir2/dir3"},
					},
					{
						dir:           "/test/dir2/dir3",
						expectedPaths: []string{"/test/dir2/dir3/dir4"},
					},
				},
			},
		})
}

func testBlobFSListFiles(t *testing.T, timeout time.Duration, location string) {

	testBlobFSList(
		t,
		"ListFiles",
		timeout,
		location,
		func(t *testing.T, f fixture.F, fs blobFS, remotepath string) []string {
			return fs.List(t, f, remotepath)
		},
		[]listTestCase{
			{
				name:        "OneFileOnRoot",
				remotefiles: []string{"/file"},
				ops: []listTestOperation{
					{
						dir:           "/",
						expectedPaths: []string{"/file"},
					},
				},
			},
			{
				name:        "OneFile",
				remotefiles: []string{"/test/file"},
				ops: []listTestOperation{
					{
						dir:           "/test",
						expectedPaths: []string{"/test/file"},
					},
				},
			},
			{
				name: "MultipleFiles",
				remotefiles: []string{
					"/test/file1",
					"/test/file2",
					"/test/file3",
				},
				ops: []listTestOperation{
					{
						dir: "/test",
						expectedPaths: []string{
							"/test/file1",
							"/test/file2",
							"/test/file3",
						},
					},
				},
			},
			{
				// WHY: The alternative is to implemented a full hierarchical FS
				// on top of azure blob, seems excessive.
				name: "PathCanBeFileAndDir",
				remotefiles: []string{
					"/test/dir",
					"/test/dir/file",
				},
				ops: []listTestOperation{
					{
						dir:           "/test",
						expectedPaths: []string{"/test/dir"},
					},
					{
						dir:           "/test/dir",
						expectedPaths: []string{"/test/dir/file"},
					},
				},
			},
			{
				name:        "EmptyDir",
				remotefiles: []string{"/test/dir"},
				ops: []listTestOperation{
					{
						dir:           "/test/dir",
						expectedPaths: []string{},
					},
				},
			},
			{
				name:        "WrongDir",
				remotefiles: []string{},
				ops: []listTestOperation{
					{
						dir:           "/",
						expectedPaths: []string{},
					},
					{
						dir:           "/test",
						expectedPaths: []string{},
					},
				},
			},
			{
				name: "NestedDirs",
				remotefiles: []string{
					"/test/file",
					"/test/dir1/file",
					"/test/dir2/file",
					"/test/dir2/dir3/file1",
					"/test/dir2/dir3/file2",
				},
				ops: []listTestOperation{
					{
						dir:           "/test",
						expectedPaths: []string{"/test/file"},
					},
					{
						dir:           "/test/dir1",
						expectedPaths: []string{"/test/dir1/file"},
					},
					{
						dir:           "/test/dir2",
						expectedPaths: []string{"/test/dir2/file"},
					},
					{
						dir: "/test/dir2/dir3",
						expectedPaths: []string{
							"/test/dir2/dir3/file1",
							"/test/dir2/dir3/file2",
						},
					},
				},
			},
		})
}

func testBlobFSUploadsWhenAccountAndContainerExists(t *testing.T, f fixture.F) {
	sf, cleanup := setupBlobStorageFixture(t, f, "Standard_LRS", "Cool")
	defer cleanup()

	expectedPath := "/test/acc/container/existent/file"
	createStorageAccountBLOB(f, sf.account, sf.sku, sf.tier)
	checkStorageBlobAccount(t, f, sf.account, sf.sku, sf.tier)
	createStorageAccountContainer(f, sf.account, sf.container)

	blobFSUpload(
		t,
		f,
		sf.account,
		"Premium_LRS",
		"Hot",
		sf.container,
		expectedPath,
		sf.testfile,
	)

	checkStorageBlobAccount(t, f, sf.account, sf.sku, sf.tier)
	filecontent := downloadFileBLOB(t, f, sf.account, sf.container, expectedPath)
	assert.EqualStrings(t, sf.testfileContent, filecontent, "checking uploaded BLOB")
}

func uploaderCreatesAccountAndContainerIfNonExistent(
	t *testing.T,
	f fixture.F,
	sku string,
	tier string,
) {
	sf, cleanup := setupBlobStorageFixture(t, f, sku, tier)
	defer cleanup()

	expectedPath := "/test/acc/container/nonexistent/file"

	blobFSUpload(
		t,
		f,
		sf.account,
		sf.sku,
		sf.tier,
		sf.container,
		expectedPath,
		sf.testfile,
	)
	checkStorageBlobAccount(t, f, sf.account, sf.sku, sf.tier)
	filecontent := downloadFileBLOB(t, f, sf.account, sf.container, expectedPath)

	assert.EqualStrings(t, sf.testfileContent, filecontent, "checking uploaded BLOB")
}

func testBlobFSCreatesAccountAndContainerIfNonExistent(
	t *testing.T,
	timeout time.Duration,
	location string,
) {
	type TestCase struct {
		sku  string
		tier string
	}

	// WHY: lot of cases are not tested because they do not work on az cli
	// like using Premium_LRS or Standard_ZRS as SKU.
	// The docs are pretty confusing right now
	// on how to use the SKU option when handling blob storage.
	tcases := []TestCase{
		TestCase{
			sku:  "Standard_LRS",
			tier: "Cool",
		},
		TestCase{
			sku:  "Standard_LRS",
			tier: "Hot",
		},
		TestCase{
			sku:  "Standard_GRS",
			tier: "Cool",
		},
		TestCase{
			sku:  "Standard_GRS",
			tier: "Hot",
		},
		TestCase{
			sku:  "Standard_RAGRS",
			tier: "Cool",
		},
		TestCase{
			sku:  "Standard_RAGRS",
			tier: "Hot",
		},
	}

	for _, tcase := range tcases {
		tname := fmt.Sprintf(
			"UploaderCreatesAccountAndContainerIfNonExistent%s%s",
			tcase.sku,
			tcase.tier,
		)
		fixture.Run(t, tname, timeout, location, func(t *testing.T, f fixture.F) {
			uploaderCreatesAccountAndContainerIfNonExistent(t, f, tcase.sku, tcase.tier)
		})
	}
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

type blobFS struct {
	account   string
	sku       string
	tier      string
	container string
}

func newBlobFS(
	account string,
	sku string,
	tier string,
	container string,
) blobFS {
	return blobFS{
		account:   account,
		sku:       sku,
		tier:      tier,
		container: container,
	}
}

func (fs *blobFS) ListDir(t *testing.T, f fixture.F, remotedir string) []string {
	res := execWithIPC(t, f, func(outputpath string) {
		f.Shell.Run(
			"./testdata/blob_fs_listdir.sh",
			f.ResGroupName,
			fs.account,
			fs.container,
			remotedir,
			outputpath,
		)
	})
	if res == "" {
		return []string{}
	}
	return strings.Split(res, " ")
}

func (fs *blobFS) List(t *testing.T, f fixture.F, remotedir string) []string {
	res := execWithIPC(t, f, func(outputpath string) {
		f.Shell.Run(
			"./testdata/blob_fs_list.sh",
			f.ResGroupName,
			fs.account,
			fs.container,
			remotedir,
			outputpath,
		)
	})
	if res == "" {
		return []string{}
	}
	return strings.Split(res, " ")
}

func (fs *blobFS) DownloadDir(
	t *testing.T,
	f fixture.F,
	localdir string,
	remotedir string,
) {
	f.Shell.Run(
		"./testdata/blob_fs_download_dir.sh",
		f.ResGroupName,
		fs.account,
		fs.container,
		localdir,
		remotedir,
	)
}

func (fs *blobFS) UploadDir(
	t *testing.T,
	f fixture.F,
	remotedir string,
	localdir string,
) {
	f.Shell.Run(
		"./testdata/blob_fs_upload_dir.sh",
		f.ResGroupName,
		f.Location,
		fs.account,
		fs.sku,
		fs.tier,
		fs.container,
		remotedir,
		localdir,
	)
}

func (fs *blobFS) Upload(
	t *testing.T,
	f fixture.F,
	remotefile string,
	localfile string,
) {
	blobFSUpload(
		t,
		f,
		fs.account,
		fs.sku,
		fs.tier,
		fs.container,
		remotefile,
		localfile,
	)
}

func blobFSUpload(
	t *testing.T,
	f fixture.F,
	account string,
	sku string,
	tier string,
	container string,
	remotepath string,
	localpath string,
) {
	f.Shell.Run(
		"./testdata/blob_fs_upload.sh",
		f.ResGroupName,
		f.Location,
		account,
		sku,
		tier,
		container,
		remotepath,
		localpath,
	)
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
