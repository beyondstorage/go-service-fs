package fs

import (
	"io/ioutil"
	"testing"

	"github.com/aos-dev/go-storage/v2/types"
	"github.com/stretchr/testify/assert"
)

func fsReaddir(b *testing.B) {
	s, _ := newStorager()

	it, err := s.ListDir("/usr/lib")
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

	it, err := s.ListDir("/usr/lib")
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
