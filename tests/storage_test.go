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

func TestCopy(t *testing.T) {
	if os.Getenv("STORAGE_FS_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_FS_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestCopier(t, setupTest(t))
	tests.TestCopierWithDir(t, setupTest(t))
}

func TestMove(t *testing.T) {
	if os.Getenv("STORAGE_FS_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_FS_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestMover(t, setupTest(t))
	tests.TestMoverWithDir(t, setupTest(t))
}

func TestLinker(t *testing.T) {
	if os.Getenv("STORAGE_FS_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_FS_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestLinker(t, setupTest(t))
}
