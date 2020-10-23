// +build aix darwin dragonfly freebsd js,wasm linux netbsd openbsd solaris

package fs

import (
	"io/ioutil"
	"os"
	"testing"
)

func fsReaddir(b *testing.B) {
	f, err := os.Open("/usr/lib")
	if err != nil {
		b.Error(err)
	}

	buf := make([]byte, 8192)

	for {
		files, err := getFiles(int(f.Fd()), buf)
		if err != nil {
			b.Error(err)
		}
		if len(files) == 0 {
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
	f, err := os.Open("/usr/lib")
	if err != nil {
		t.Error(err)
	}

	buf := make([]byte, 8192)

	for {
		files, err := getFiles(int(f.Fd()), buf)
		if err != nil {
			t.Error(err)
		}
		if len(files) == 0 {
			break
		}
		for _, v := range files {
			t.Logf("%v", v)
		}
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
