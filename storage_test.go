package fs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/aos-dev/go-storage/v2/pairs"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	typ "github.com/aos-dev/go-storage/v2/types"
)

func TestStorage_String(t *testing.T) {
	c := Storage{}
	c.workDir = "/test"

	assert.Equal(t, "Storager fs {WorkDir: /test}", c.String())
}

func TestStorage_Metadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	{
		client := Storage{workDir: "/test"}

		m, err := client.Metadata()
		assert.NoError(t, err)
		assert.Equal(t, "/test", m.WorkDir)
	}
}

type fileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (f fileInfo) Name() string {
	return f.name
}

func (f fileInfo) Size() int64 {
	return f.size
}

func (f fileInfo) Mode() os.FileMode {
	return f.mode
}

func (f fileInfo) ModTime() time.Time {
	return f.modTime
}

func (f fileInfo) IsDir() bool {
	return f.mode.IsDir()
}

func (f fileInfo) Sys() interface{} {
	return f
}

func TestStorage_Stat(t *testing.T) {
	nowTime := time.Now()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string
		err  error
		file fileInfo

		object *typ.Object
	}{
		{
			"regular file",
			nil,
			fileInfo{
				name:    "regular file",
				size:    1234,
				mode:    0777,
				modTime: nowTime,
			},
			typ.NewObject(nil, true).
				SetID("regular file").
				SetName("regular file").
				SetType(typ.ObjectTypeFile).
				SetContentType("application/octet-stream").
				SetSize(1234).
				SetUpdatedAt(nowTime),
		},
		{
			"dir",
			nil,
			fileInfo{
				name:    "dir",
				size:    0,
				mode:    os.ModeDir | 0777,
				modTime: nowTime,
			},
			&typ.Object{
				ID:   "dir",
				Name: "dir",
				Type: typ.ObjectTypeDir,
			},
		},
		{
			"stream",
			nil,
			fileInfo{
				name:    "stream",
				size:    0,
				mode:    os.ModeDevice | 0777,
				modTime: nowTime,
			},
			&typ.Object{
				ID:   "stream",
				Name: "stream",
				Type: typ.ObjectTypeStream,
			},
		},
		{
			"-",
			nil,
			fileInfo{},
			&typ.Object{
				ID:   "-",
				Name: "-",
				Type: typ.ObjectTypeStream,
			},
		},
		{
			"invalid",
			nil,
			fileInfo{
				name:    "invalid",
				size:    0,
				mode:    os.ModeIrregular | 0777,
				modTime: nowTime,
			},
			&typ.Object{
				ID:   "invalid",
				Name: "invalid",
				Type: typ.ObjectTypeInvalid,
			},
		},
		{
			"error",
			&os.PathError{
				Op:   "stat",
				Path: "invalid",
				Err:  os.ErrPermission,
			},
			fileInfo{}, nil,
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			client := Storage{
				osStat: func(name string) (os.FileInfo, error) {
					assert.Equal(t, v.name, name)
					return v.file, v.err
				},
			}
			o, err := client.Stat(v.name)
			assert.Equal(t, v.err == nil, err == nil)
			if v.object != nil {
				assert.NotNil(t, o)
				// FIXME: we need to have a trick to test values.
				// assert.EqualValues(t, v.object, o)
			} else {
				assert.Nil(t, o)
			}
		})
	}
}

func TestStorage_WriteStream(t *testing.T) {
	err := os.Remove("/tmp/test")
	var e *os.PathError
	if errors.As(err, &e) {
		t.Logf("%#v", e)
	}
}

func TestStorage_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string
		err  error
	}{
		{"delete file", nil},
		{"delete nonempty dir", &os.PathError{
			Op:   "remove",
			Path: "delete nonempty dir",
			Err:  errors.New("remove fail"),
		}},
	}

	for _, v := range tests {
		v := v

		t.Run(v.name, func(t *testing.T) {

			client := Storage{
				osRemove: func(name string) error {
					assert.Equal(t, v.name, name)
					return v.err
				},
			}
			err := client.Delete(v.name)
			assert.Equal(t, v.err == nil, err == nil)
		})
	}
}

func TestStorage_Copy(t *testing.T) {
	t.Run("Failed at open source file", func(t *testing.T) {
		srcName := uuid.New().String()
		dstName := uuid.New().String()
		client := Storage{
			osOpen: func(name string) (file *os.File, e error) {
				assert.Equal(t, srcName, name)
				return nil, &os.PathError{
					Op:  "open",
					Err: errors.New("path error"),
				}
			},
			osMkdirAll: func(path string, perm os.FileMode) error {
				return nil
			},
		}

		err := client.Copy(srcName, dstName)
		assert.Error(t, err)
	})

	t.Run("Failed at open dst file", func(t *testing.T) {
		srcName := uuid.New().String()
		dstName := uuid.New().String()
		client := Storage{
			osOpen: func(name string) (file *os.File, e error) {
				assert.Equal(t, srcName, name)
				return nil, nil
			},
			osCreate: func(name string) (file *os.File, e error) {
				assert.Equal(t, dstName, name)
				return nil, &os.PathError{
					Op:  "open",
					Err: errors.New("open fail"),
				}
			},
			osMkdirAll: func(path string, perm os.FileMode) error {
				return nil
			},
		}

		err := client.Copy(srcName, dstName)
		assert.Error(t, err)
	})

	t.Run("Failed at io.CopyBuffer", func(t *testing.T) {
		srcName := uuid.New().String()
		dstName := uuid.New().String()
		client := Storage{
			osOpen: func(name string) (file *os.File, e error) {
				assert.Equal(t, srcName, name)
				return nil, nil
			},
			osCreate: func(name string) (file *os.File, e error) {
				assert.Equal(t, dstName, name)
				return nil, nil
			},
			ioCopyBuffer: func(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
				return 0, io.ErrShortWrite
			},
			osMkdirAll: func(path string, perm os.FileMode) error {
				return nil
			},
		}

		err := client.Copy(srcName, dstName)
		assert.Error(t, err)
	})

	t.Run("All successful", func(t *testing.T) {
		fakeFile := &os.File{}
		// Monkey patch the file's Close.
		monkey.PatchInstanceMethod(reflect.TypeOf(fakeFile), "Close",
			func(f *os.File) error {
				return nil
			})

		srcName := uuid.New().String()
		dstName := uuid.New().String()
		client := Storage{
			osOpen: func(name string) (file *os.File, e error) {
				assert.Equal(t, srcName, name)
				return fakeFile, nil
			},
			osCreate: func(name string) (file *os.File, e error) {
				assert.Equal(t, dstName, name)
				return fakeFile, nil
			},
			ioCopyBuffer: func(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
				return 0, nil
			},
			osMkdirAll: func(path string, perm os.FileMode) error {
				return nil
			},
		}

		err := client.Copy(srcName, dstName)
		assert.NoError(t, err)
	})
}

func TestStorage_Move(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		srcName := uuid.New().String()
		dstName := uuid.New().String()

		client := Storage{
			osRename: func(oldpath, newpath string) error {
				assert.Equal(t, srcName, oldpath)
				assert.Equal(t, dstName, newpath)
				return &os.LinkError{
					Op:  "rename",
					Old: oldpath,
					New: newpath,
					Err: errors.New("rename fail"),
				}
			},
			osMkdirAll: func(path string, perm os.FileMode) error {
				return nil
			},
		}

		err := client.Move(srcName, dstName)
		assert.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		srcName := uuid.New().String()
		dstName := uuid.New().String()

		client := Storage{
			osRename: func(oldpath, newpath string) error {
				assert.Equal(t, srcName, oldpath)
				assert.Equal(t, dstName, newpath)
				return nil
			},
			osMkdirAll: func(path string, perm os.FileMode) error {
				return nil
			},
		}

		err := client.Move(srcName, dstName)
		assert.NoError(t, err)
	})
}

func TestStorage_ListDir(t *testing.T) {
	paths := make([]string, 100)
	for k := range paths {
		paths[k] = uuid.New().String()
	}

	tests := []struct {
		name             string
		enableFollowLink bool
		fi               []os.FileInfo
		items            []*typ.Object
		err              error
	}{
		{
			"success file",
			false,
			[]os.FileInfo{
				fileInfo{
					name:    "test_file",
					size:    1234,
					mode:    0644,
					modTime: time.Unix(1, 0),
				},
			},
			[]*typ.Object{
				typ.NewObject(nil, true).
					SetID(filepath.Join(paths[0], "test_file")).
					SetName(path.Join(paths[0], "test_file")).
					SetType(typ.ObjectTypeFile).
					SetContentType("application/octet-stream").
					SetSize(1234).
					SetUpdatedAt(time.Unix(1, 0)),
			},
			nil,
		},
		{
			"success file recursively",
			false,
			[]os.FileInfo{
				fileInfo{
					name:    "test_file",
					size:    1234,
					mode:    0644,
					modTime: time.Unix(1, 0),
				},
			},
			[]*typ.Object{
				typ.NewObject(nil, true).
					SetID(filepath.Join(paths[1], "test_file")).
					SetName(path.Join(paths[1], "test_file")).
					SetType(typ.ObjectTypeFile).
					SetContentType("application/octet-stream").
					SetSize(1234).
					SetUpdatedAt(time.Unix(1, 0)),
			},
			nil,
		},
		{
			"success dir",
			false,
			[]os.FileInfo{
				fileInfo{
					name:    "test_dir",
					size:    0,
					mode:    os.ModeDir | 0755,
					modTime: time.Unix(1, 0),
				},
			},
			[]*typ.Object{
				{
					ID:   filepath.Join(paths[2], "test_dir"),
					Name: path.Join(paths[2], "test_dir"),
					Type: typ.ObjectTypeDir,
				},
			},
			nil,
		},
		{
			"success dir recursively",
			false,
			[]os.FileInfo{
				fileInfo{
					name:    "test_dir",
					size:    0,
					mode:    os.ModeDir | 0755,
					modTime: time.Unix(1, 0),
				},
			},
			[]*typ.Object{
				{
					ID:   filepath.Join(paths[3], "test_dir"),
					Name: path.Join(paths[3], "test_dir"),
					Type: typ.ObjectTypeDir,
				},
			},
			nil,
		},
		{
			"success file under windows",
			false,
			[]os.FileInfo{
				fileInfo{
					name:    "test_file",
					size:    1234,
					mode:    0644,
					modTime: time.Unix(1, 0),
				},
			},
			[]*typ.Object{
				typ.NewObject(nil, true).
					SetID(filepath.Join(paths[4], "test_file")).
					// Make sure ListDir return a name with slash.
					SetName(fmt.Sprintf("%s/%s", paths[4], "test_file")).
					SetType(typ.ObjectTypeFile).
					SetContentType("application/octet-stream").
					SetSize(1234).
					SetUpdatedAt(time.Unix(1, 0)),
			},
			nil,
		},
		{
			"skip link",
			false,
			[]os.FileInfo{
				fileInfo{
					name:    "test_link",
					size:    0,
					mode:    os.ModeSymlink,
					modTime: time.Unix(1, 0),
				},
			},
			[]*typ.Object{},
			nil,
		},
		{
			"follow link",
			true,
			[]os.FileInfo{
				fileInfo{
					name:    "test_link",
					size:    1234,
					mode:    os.ModeSymlink,
					modTime: time.Unix(1, 0),
				},
			},
			[]*typ.Object{
				typ.NewObject(nil, true).
					SetID(filepath.Join(paths[6], "test_link")).
					// Make sure ListDir return a name with slash.
					SetName(fmt.Sprintf("%s/%s", paths[6], "test_link")).
					SetType(typ.ObjectTypeFile).
					SetContentType("application/octet-stream").
					SetSize(1234).
					SetUpdatedAt(time.Unix(1, 0)),
			},
			nil,
		},
		{
			"os error",
			false,
			nil,
			[]*typ.Object{},
			&os.PathError{Op: "readdir", Path: "", Err: errors.New("readdir fail")},
		},
	}

	monkey.Patch(filepath.EvalSymlinks, func(s string) (string, error) {
		return s, nil
	})
	defer monkey.UnpatchAll()

	for k, v := range tests {
		monkey.Patch(os.Stat, func(s string) (os.FileInfo, error) {
			for _, o := range v.fi {
				if strings.HasSuffix(s, o.Name()) {
					return o, nil
				}
			}
			return nil, os.ErrNotExist
		})
		t.Run(v.name, func(t *testing.T) {
			client := &Storage{
				ioutilReadDir: func(dirname string) (infos []os.FileInfo, e error) {
					assert.Equal(t, paths[k], dirname)
					return v.fi, v.err
				},
			}

			items := make([]*typ.Object, 0)

			it, err := client.ListDir(paths[k], WithEnableLinkFollow(v.enableFollowLink))
			if err != nil {
				t.Error(err)
			}

			for {
				o, err := it.Next()
				if err == typ.IterateDone {
					break
				}
				assert.Equal(t, v.err == nil, err == nil)
				if err != nil {
					break
				}

				items = append(items, o)
			}
			// FIXME: we need test values here
			// assert.EqualValues(t, v.items, items)
		})
	}
}

func TestStorage_Read(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		pairs   []*typ.Pair
		isNil   bool
		openErr error
		seekErr error
	}{
		{
			"success",
			"test_success",
			nil,
			false,
			nil,
			nil,
		},
		{
			"error",
			"test_error",
			nil,
			true,
			&os.PathError{Op: "readdir", Path: "", Err: errors.New("readdir fail")},
			nil,
		},
		{
			"stdin",
			"-",
			nil,
			false,
			nil,
			nil,
		},
		{
			"stdin with size",
			"-",
			[]*typ.Pair{
				pairs.WithSize(100),
			},
			false,
			nil,
			nil,
		},
		{
			"success with size",
			"test_success",
			[]*typ.Pair{
				pairs.WithSize(100),
			},
			false,
			nil,
			nil,
		},
		{
			"success with offset",
			"test_success",
			[]*typ.Pair{
				pairs.WithOffset(10),
			},
			false,
			nil,
			nil,
		},
		{
			"error with offset",
			"test_success",
			[]*typ.Pair{
				pairs.WithOffset(10),
			},
			true,
			nil,
			io.ErrUnexpectedEOF,
		},
		{
			"success with and size offset",
			"test_success",
			[]*typ.Pair{
				pairs.WithSize(100),
				pairs.WithOffset(10),
			},
			false,
			nil,
			nil,
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			fakeFile := &os.File{}
			monkey.PatchInstanceMethod(reflect.TypeOf(fakeFile), "Seek", func(f *os.File, offset int64, whence int) (ret int64, err error) {
				t.Logf("Seek has been called.")
				assert.Equal(t, int64(10), offset)
				assert.Equal(t, 0, whence)
				return 0, v.seekErr
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(fakeFile), "Read", func(f *os.File, b []byte) (n int, err error) {
				t.Logf("Read has been called.")
				b = append(b, []byte("xxxx")...)
				return 4, io.EOF
			})

			client := Storage{
				osOpen: func(name string) (file *os.File, e error) {
					assert.Equal(t, v.path, name)
					return fakeFile, v.openErr
				},
			}

			var buf bytes.Buffer
			n, err := client.Read(v.path, &buf, v.pairs...)
			assert.Equal(t, v.openErr == nil && v.seekErr == nil, err == nil)
			assert.Equal(t, int64(buf.Len()), n)
		})
	}
}

func TestStorage_Write(t *testing.T) {
	paths := make([]string, 10)
	for k := range paths {
		paths[k] = uuid.New().String()
	}

	tests := []struct {
		name         string
		osCreate     func(name string) (*os.File, error)
		ioCopyN      func(dst io.Writer, src io.Reader, n int64) (written int64, err error)
		ioCopyBuffer func(dst io.Writer, src io.Reader, buf []byte) (written int64, err error)
		hasErr       bool
		written      int64
	}{
		{
			"failed os create",
			func(name string) (file *os.File, e error) {
				assert.Equal(t, paths[0], name)
				return nil, &os.PathError{
					Op:   "open",
					Path: "",
					Err:  os.ErrNotExist,
				}
			},
			nil,
			nil,
			true,
			0,
		},
		{
			"failed io copyn",
			func(name string) (file *os.File, e error) {
				assert.Equal(t, paths[1], name)
				return &os.File{}, nil
			},
			func(dst io.Writer, src io.Reader, n int64) (written int64, err error) {
				return 0, io.EOF
			},
			nil,
			true,
			0,
		},
		{
			"failed io copy buffer",
			nil,
			nil,
			func(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
				return 0, io.EOF
			},
			true,
			0,
		},
		{
			"success with size",
			func(name string) (file *os.File, e error) {
				assert.Equal(t, paths[3], name)
				return &os.File{}, nil
			},
			func(dst io.Writer, src io.Reader, n int64) (written int64, err error) {
				assert.Equal(t, int64(1234), n)
				return n, nil
			},
			nil,
			false,
			1234,
		},
		{
			"success with stdout",
			nil,
			func(dst io.Writer, src io.Reader, n int64) (written int64, err error) {
				assert.Equal(t, int64(1234), n)
				return n, nil
			},
			nil,
			false,
			1234,
		},
	}

	for k, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			client := Storage{
				osCreate:     v.osCreate,
				ioCopyN:      v.ioCopyN,
				ioCopyBuffer: v.ioCopyBuffer,
				osMkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
			}

			var pair []*typ.Pair
			if v.ioCopyN != nil {
				pair = append(pair, pairs.WithSize(1234))
			}

			var err error
			var n int64
			if v.osCreate == nil {
				n, err = client.Write("-", nil, pair...)
			} else {
				n, err = client.Write(paths[k], nil, pair...)
			}
			assert.Equal(t, v.hasErr, err != nil)
			assert.Equal(t, v.written, n)
		})
	}
}
