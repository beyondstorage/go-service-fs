// +build integration_test

package tests

import (
	"testing"

	tests "github.com/aos-dev/go-integration-test/v2"
)

func TestStorage(t *testing.T) {
	tests.TestStorager(t, setupTest(t))
}
