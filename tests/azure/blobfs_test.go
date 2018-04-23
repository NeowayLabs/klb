package azure_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/NeowayLabs/klb/tests/lib/assert"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func TestBlobFS(t *testing.T) {
	timeout := 60 * time.Minute
	t.Parallel()
	fixture.Run(
		t,
		"UploadsWhenAccountAndContainerExists",
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

func testBlobFSUploadsWhenAccountAndContainerExists(t *testing.T, f fixture.F) {
	sf, cleanup := setupBlobStorageFixture(t, f, "Standard_LRS", "Cool")
	defer cleanup()

	expectedPath := "/test/acc/container/existent/file"
	createStorageAccountBLOB(f, sf.account, sf.sku, sf.tier)
	checkStorageBlobAccount(t, f, sf.account, sf.sku, sf.tier)
	createStorageAccountContainer(f, sf.account, sf.container)

	fs := newBlobFS(sf.account, "Premium_LRS", "Hot", sf.container)
	fs.Upload(
		t,
		f,
		expectedPath,
		sf.testfile,
	)

	checkStorageBlobAccount(t, f, sf.account, sf.sku, sf.tier)
	filecontent := fs.Download(t, f, expectedPath)
	assert.EqualStrings(t, sf.testfileContent, filecontent, "checking uploaded BLOB")
}

func blobFSCreatesAccountAndContainerIfNonExistent(
	t *testing.T,
	f fixture.F,
	sku string,
	tier string,
) {
	sf, cleanup := setupBlobStorageFixture(t, f, sku, tier)
	defer cleanup()

	expectedPath := "/test/acc/container/nonexistent/file"
	fs := newBlobFS(sf.account, sf.sku, sf.tier, sf.container)

	fs.Upload(t, f, expectedPath, sf.testfile)
	checkStorageBlobAccount(t, f, sf.account, sf.sku, sf.tier)
	filecontent := fs.Download(t, f, expectedPath)

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
			"CreatesAccountAndContainerIfNonExistent%s%s",
			tcase.sku,
			tcase.tier,
		)
		fixture.Run(t, tname, timeout, location, func(t *testing.T, f fixture.F) {
			blobFSCreatesAccountAndContainerIfNonExistent(t, f, tcase.sku, tcase.tier)
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
				name: "MultipleFilesOnRoot",
				remotefiles: []string{
					"/file",
					"/file2",
					"/file3",
				},
				ops: []listTestOperation{
					{
						dir: "/",
						expectedPaths: []string{
							"/file",
							"/file2",
							"/file3",
						},
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
	container := fixture.NewUniqueName("testcontainer")
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
		tname := testprefix + test.name
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

func (fs *blobFS) Download(
	t *testing.T,
	f fixture.F,
	remotefile string,
) string {
	return execWithIPC(t, f, func(outputpath string) {
		f.Shell.Run(
			"./testdata/blob_fs_download.sh",
			f.ResGroupName,
			fs.account,
			fs.container,
			remotefile,
			outputpath,
		)
	})
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

	tests := []TestCase{
		{
			name: "OneFile",
			remoteFiles: []string{
				"/one/test1",
			},
			downloadDir: "/one",
			wantedFiles: []string{
				"/test1",
			},
		},
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
		{
			name: "Level1",
			remoteFiles: []string{
				"/test1",
				"/test2",
				"/dir/test1",
				"/dir/test2",
			},
			downloadDir: "/dir",
			wantedFiles: []string{
				"/test1",
				"/test2",
			},
		},
		{
			name: "Level2",
			remoteFiles: []string{
				"/test1",
				"/test2",
				"/dir/test1",
				"/dir/test2",
				"/dir/dir2/test1",
				"/dir/dir2/test2",
			},
			downloadDir: "/dir/dir2",
			wantedFiles: []string{
				"/test1",
				"/test2",
			},
		},
		{
			name: "SubdirsAreDownloaded",
			remoteFiles: []string{
				"/test1",
				"/test2",
				"/dir/test1",
				"/dir/test2",
				"/dir/dir2/test1",
				"/dir/dir2/test2",
			},
			downloadDir: "/dir",
			wantedFiles: []string{
				"/test1",
				"/test2",
				"/dir2/test1",
				"/dir2/test2",
			},
		},
	}

	for _, test := range tests {
		name := "DownloadDir" + test.name
		remoteFiles := test.remoteFiles
		wantedFiles := test.wantedFiles
		downloadDir := test.downloadDir

		fixture.Run(t, name, timeout, location, func(t *testing.T, f fixture.F) {
			fs := setupFS()

			createStorageAccountBLOB(f, fs.account, fs.sku, fs.tier)
			checkStorageBlobAccount(t, f, fs.account, fs.sku, fs.tier)
			createStorageAccountContainer(f, fs.account, fs.container)

			for _, remoteFile := range remoteFiles {
				localFile, cleanup := setupTestFile(t, remoteFile)
				fs.Upload(t, f, remoteFile, localFile)
				cleanup()
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
				gotFiles = append(gotFiles, path)
				return nil
			})
			assert.NoError(t, err)

			if len(wantedFiles) != len(gotFiles) {
				t.Errorf("wanted files and got files len don't match")
				t.Fatalf("want files[%s] got[%s]", wantedFiles, gotFiles)
			}
			
			removeFilePrefix := func(path string) string {
				return strings.TrimPrefix(path, tmpdir)
			}
			
			removeFilesPrefix := func(paths []string) []string {
				trimmed := make([]string, len(paths))
				for i, path := range paths {
					trimmed[i] = removeFilePrefix(path)
				}
				return trimmed
			}
			
			expectedContent := func(path string) string {
				// WHY: we use the original remote file path as
				// the content itself (this way the content is unique for each file).
				return filepath.Join(downloadDir, path)
			}

			for _, wantedFile := range wantedFiles {

				got := false
				for _, gotFile := range gotFiles {

					gotFileRelative := removeFilePrefix(gotFile)
					if gotFileRelative != wantedFile {
						continue
					}
					
					got = true
					f, err := os.Open(gotFile)
					assert.NoError(t, err, "opening file")
					defer f.Close()
						
					contentsr, err := ioutil.ReadAll(f)
					assert.NoError(t, err, "reading file")
						
					contents := string(contentsr)
					wantedContents := expectedContent(wantedFile)
						
					if contents != wantedContents {
						t.Fatalf(
							"wanted file[%s] contents[%s] got[%s]",
							wantedFile,
							wantedContents,
							contents,
						)
					}
					
				}
				if !got {
					t.Fatalf(
						"wanted files[%s] != got[%s]",
						wantedFiles,
						removeFilesPrefix(gotFiles),
					)
				}
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
		localfiles []string
		uploadDir string
		wantedRemoteFiles []string
	}

	tests := []TestCase{
		{
			name: "RegressionFilesWithDotAndSlashes",
			uploadDir: "/regression/slashes/dots",
			localfiles: []string{
				"/test/1",
				"/test/model.ckpt.data-00000-of-00001",
			},
			wantedRemoteFiles: []string{
				"/regression/slashes/dots/test/1",
				"/regression/slashes/dots/test/model.ckpt.data-00000-of-00001",
			},
		},
		{
			name: "OneDirToRoot",
			uploadDir: "/",
			localfiles: []string{
				"/test/1",
				"/test/2",
			},
			wantedRemoteFiles: []string{
				"/test/1",
				"/test/2",
			},
		},
		{
			name: "MultiplesDirsToRoot",
			uploadDir: "/",
			localfiles: []string{
				"/test/1",
				"/test/2",
				"/test2/1",
				"/test2/2",
			},
			wantedRemoteFiles: []string{
				"/test/1",
				"/test/2",
				"/test2/1",
				"/test2/2",
			},
		},
		{
			name: "MultipleDirsNested",
			uploadDir: "/remotedir/nested",
			localfiles: []string{
				"/test/nested/1",
				"/test/nested/2",
				"/test/nested2/1",
				"/other/1",
				"/other/nest1/nest2/1",
			},
			wantedRemoteFiles: []string{
				"/remotedir/nested/test/nested/1",
				"/remotedir/nested/test/nested/2",
				"/remotedir/nested/test/nested2/1",
				"/remotedir/nested/other/1",
				"/remotedir/nested/other/nest1/nest2/1",
			},
		},
	}

	for _, tt := range tests {
		test := tt
		name := "UploadDir" + test.name
		fixture.Run(t, name, timeout, location, func(t *testing.T, f fixture.F) {
			tdir, err := ioutil.TempDir("", "uploadDirTests")
			assert.NoError(t, err)
			defer func() {
				assert.NoError(t, os.RemoveAll(tdir), "removing tmp dir")
			}()
			
			for _, localfile := range test.localfiles {
				content := fixture.NewUniqueName("random-content")
				fullpath := filepath.Join(tdir, localfile)
				basedir := filepath.Dir(fullpath)
				assert.NoError(t, os.MkdirAll(basedir, 0644), "creating dir for local file")
				assert.NoError(t, ioutil.WriteFile(
						fullpath,
						[]byte(content),
						0644,
					),
					"writing local test file",
				)
			}
			
			fs := setupFS()
			fs.UploadDir(t, f, test.uploadDir, tdir)
			
			getEquivalentLocalFile := func(remotefile string) string {
				return filepath.Join(tdir, strings.TrimPrefix(remotefile, test.uploadDir))
			}
			
			for _, wantedRemoteFile := range test.wantedRemoteFiles {
				equivalentLocalFile := getEquivalentLocalFile(wantedRemoteFile)
				wantedContentRaw, err := ioutil.ReadFile(equivalentLocalFile)
				assert.NoError(t, err, "reading local file")
				
				wantedContent := string(wantedContentRaw)
				
				gotContent := fs.Download(t, f, wantedRemoteFile)
				
				if wantedContent != gotContent {
					t.Fatalf(
						"remote file[%s] contents[%s] != local file[%s] contents[%s]",
						wantedRemoteFile,
						gotContent,
						equivalentLocalFile,
						wantedContent,
					)
				}
				
			}
		})
	}

}
