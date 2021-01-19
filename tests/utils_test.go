// +build integration_test

package tests

import (
	"testing"

	fs "github.com/aos-dev/go-service-fs"
	ps "github.com/aos-dev/go-storage/v3/pairs"
	"github.com/aos-dev/go-storage/v3/types"
)

func setupTest(t *testing.T) types.Storager {
	tmpDir := t.TempDir()
	t.Logf("Setup test at %s", tmpDir)

	store, err := fs.NewStorager(ps.WithWorkDir(tmpDir))
	if err != nil {
		t.Errorf("new storager: %v", err)
	}
	return store
}
