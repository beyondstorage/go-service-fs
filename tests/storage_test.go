package tests

import (
	"os"
	"testing"

	tests "github.com/beyondstorage/go-integration-test/v4"
)

func TestStorage(t *testing.T) {
	if os.Getenv("STORAGE_FS_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_FS_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestStorager(t, setupTest(t))
}

func TestAppend(t *testing.T) {
	if os.Getenv("STORAGE_FS_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_FS_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestAppender(t, setupTest(t))
}

func TestDir(t *testing.T) {
	if os.Getenv("STORAGE_FS_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_FS_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestDirer(t, setupTest(t))
}
