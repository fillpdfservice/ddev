package archive_test

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/drud/ddev/pkg/archive"
	"github.com/drud/ddev/pkg/testcommon"
	"github.com/drud/ddev/pkg/util"
	"github.com/stretchr/testify/assert"
)

var (
	// TestTarArchiveURL provides the URL of the test tar.gz asset
	TestTarArchiveURL = "https://github.com/drud/wordpress/archive/v0.4.0.tar.gz"
	// TestTarArchivePath provides the path the test tar.gz asset is downloaded to
	TestTarArchivePath string
	// TestTarArchiveExtractDir is the directory in the archive to extract
	TestTarArchiveExtractDir = "wordpress-0.4.0/"
)

func TestMain(m *testing.M) {
	// prep assets for files tests
	testPath, err := ioutil.TempDir("", "filetest")
	util.CheckErr(err)
	testPath, err = filepath.EvalSymlinks(testPath)
	util.CheckErr(err)
	testPath = filepath.Clean(testPath)
	TestTarArchivePath = filepath.Join(testPath, "files.tar.gz")

	err = util.DownloadFile(TestTarArchivePath, TestTarArchiveURL)
	if err != nil {
		log.Fatalf("archive download failed: %s", err)
	}

	testRun := m.Run()

	// cleanup test file
	err = os.Remove(TestTarArchivePath)
	if err != nil {
		log.Fatal("failed to remove test asset: ", err)
	}

	os.Exit(testRun)
}

// TestUntar tests untar functionality, including the starting directory
func TestUntar(t *testing.T) {
	assert := assert.New(t)
	exDir := testcommon.CreateTmpDir("TestUnTar1")

	err := archive.Untar(TestTarArchivePath, exDir, "")
	assert.NoError(err)

	// Make sure that our base extraction directory is there
	finfo, err := os.Stat(filepath.Join(exDir, TestTarArchiveExtractDir))
	assert.NoError(err)
	assert.True(err == nil && finfo.IsDir())
	finfo, err = os.Stat(filepath.Join(exDir, TestTarArchiveExtractDir, ".ddev/config.yaml"))
	assert.NoError(err)
	assert.True(err == nil && !finfo.IsDir())

	err = os.RemoveAll(exDir)
	assert.NoError(err)

	// Now do the untar with an extraction root
	exDir = testcommon.CreateTmpDir("TestUnTar2")
	err = archive.Untar(TestTarArchivePath, exDir, TestTarArchiveExtractDir)
	assert.NoError(err)

	finfo, err = os.Stat(filepath.Join(exDir, ".ddev"))
	assert.NoError(err)
	assert.True(err == nil && finfo.IsDir())
	finfo, err = os.Stat(filepath.Join(exDir, ".ddev/config.yaml"))
	assert.NoError(err)
	assert.True(err == nil && !finfo.IsDir())

	err = os.RemoveAll(exDir)
	assert.NoError(err)

}

// TestUnzip tests unzip functionality, including the starting extraction-skip directory
func TestUnzip(t *testing.T) {

	// testZipExtractDir is the directory we may want to use to start extracting.
	testZipExtractDir := "dir2/"

	assert := assert.New(t)
	exDir := testcommon.CreateTmpDir("TestUnzip1")

	zipfilePath := filepath.Join("test", "testfile.zip")

	err := archive.Unzip(zipfilePath, exDir, "")
	assert.NoError(err)

	// Make sure that our base extraction directory is there
	finfo, err := os.Stat(filepath.Join(exDir, testZipExtractDir))
	assert.NoError(err)
	assert.True(err == nil && finfo.IsDir())
	finfo, err = os.Stat(filepath.Join(exDir, testZipExtractDir, "dir2_file.txt"))
	assert.NoError(err)
	assert.True(err == nil && !finfo.IsDir())

	err = os.RemoveAll(exDir)
	assert.NoError(err)

	// Now do the unzip with an extraction root
	exDir = testcommon.CreateTmpDir("TestUnzip2")
	err = archive.Unzip(zipfilePath, exDir, testZipExtractDir)
	assert.NoError(err)

	// Only the dir2_file should remain
	finfo, err = os.Stat(filepath.Join(exDir, "dir2_file.txt"))
	assert.NoError(err)
	assert.True(err == nil && !finfo.IsDir())

	err = os.RemoveAll(exDir)
	assert.NoError(err)
}