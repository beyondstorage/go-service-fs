package fs

import (
	ps "github.com/beyondstorage/go-storage/v4/pairs"
	"github.com/google/uuid"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"testing"

	"github.com/beyondstorage/go-storage/v4/types"
	"github.com/stretchr/testify/assert"
)

func fsReaddir(b *testing.B) {
	s, _ := newStorager()

	it, err := s.List("/usr/lib")
	if err != nil {
		b.Error(err)
	}

	for {
		_, err := it.Next()
		if err == types.IterateDone {
			break
		}
	}
}

func osReaddir(b *testing.B) {
	_, err := ioutil.ReadDir("/usr/lib")
	if err != nil {
		b.Error(err)
	}
}

func TestGetFilesFs(t *testing.T) {
	s, _ := newStorager()

	it, err := s.List("/usr/lib")
	if err != nil {
		t.Error(err)
	}

	for {
		o, err := it.Next()
		if err == types.IterateDone {
			break
		}
		assert.NotNil(t, o)
	}
}

// This test case intends to reproduce issue #68.
//
// ref: https://github.com/beyondstorage/go-service-fs/issues/68
func TestIssue68(t *testing.T) {
	tmpDir := t.TempDir()

	store, err := newStorager(ps.WithWorkDir(tmpDir))
	if err != nil {
		t.Errorf("new storager: %v", err)
	}

	// We will create upto 1000 files, introduce rand for fuzzing.
	numbers := 225 + rand.Intn(800)

	t.Logf("this tes case list %d files", numbers)

	// Create enough files in a dir, the file name must be long enough.
	// So that the total size will bigger than 8196.
	for i := 0; i < numbers; i++ {
		// uuid's max size is 36.
		// We use rand here for fuzzing.
		filename := uuid.NewString()[:1+rand.Intn(35)]

		f, err := os.Create(path.Join(tmpDir, filename))
		if err != nil {
			t.Error(err)
		}
		err = f.Close()
		if err != nil {
			t.Error(err)
		}
	}

	expected := make(map[string]struct{})
	fi, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		t.Error(err)
	}
	for _, v := range fi {
		expected[v.Name()] = struct{}{}
	}

	actual := make(map[string]struct{})
	it, err := store.List("")
	if err != nil {
		t.Error(err)
	}
	for {
		o, err := it.Next()
		if err == types.IterateDone {
			break
		}
		_, exist := actual[o.Path]
		if exist {
			t.Errorf("file %s exists already", o.Path)
			return
		}

		actual[o.Path] = struct{}{}
	}

	assert.Equal(t, expected, actual)
}

func BenchmarkGetFilesFs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fsReaddir(b)
	}
}

func BenchmarkGetFilesOs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		osReaddir(b)
	}
}
