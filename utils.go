package fs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aos-dev/go-storage/v2/services"
	typ "github.com/aos-dev/go-storage/v2/types"
)

// Std{in/out/err} support
const (
	Stdin  = "/dev/stdin"
	Stdout = "/dev/stdout"
	Stderr = "/dev/stderr"
)

// Storage is the fs client.
type Storage struct {
	// options for this storager.
	workDir string // workDir dir for all operation.
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf("Storager fs {WorkDir: %s}", s.workDir)
}

// NewStorager will create Storager only.
func NewStorager(pairs ...typ.Pair) (typ.Storager, error) {
	return newStorager(pairs...)
}

// newStorager will create a fs client.
func newStorager(pairs ...typ.Pair) (store *Storage, err error) {
	defer func() {
		if err != nil {
			err = &services.InitError{Op: "new_storager", Type: Type, Err: err, Pairs: pairs}
		}
	}()
	opt, err := parsePairStorageNew(pairs)
	if err != nil {
		return
	}

	store = &Storage{
		workDir: "/",
	}

	if opt.HasWorkDir {
		store.workDir = opt.WorkDir
	}

	// Check and create work dir
	err = os.MkdirAll(store.workDir, 0755)
	if err != nil {
		return nil, err
	}
	return
}

func formatError(err error) error {
	// Handle error returned by os package.
	switch {
	case os.IsNotExist(err):
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	case os.IsPermission(err):
		return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
	default:
		return err
	}
}

func (s *Storage) newObject(done bool) *typ.Object {
	return typ.NewObject(s, done)
}

func (s *Storage) openFile(absPath string) (f *os.File, err error) {
	switch absPath {
	case Stdin:
		f = os.Stdin
	case Stdout:
		f = os.Stdout
	case Stderr:
		f = os.Stderr
	default:
		f, err = os.Open(absPath)
	}

	return
}

func (s *Storage) createFile(absPath string) (f *os.File, err error) {
	switch absPath {
	case Stdin:
		f = os.Stdin
	case Stdout:
		f = os.Stdout
	case Stderr:
		f = os.Stderr
	default:
		defer func() {
			err = s.formatError("create_file", err, absPath)
		}()

		// Create dir before create file
		err = os.MkdirAll(filepath.Dir(absPath), 0755)
		if err != nil {
			return nil, err
		}

		f, err = os.Create(absPath)
	}

	return
}

func (s *Storage) statFile(absPath string) (fi os.FileInfo, err error) {
	switch absPath {
	case Stdin:
		fi, err = os.Stdin.Stat()
	case Stdout:
		fi, err = os.Stdout.Stat()
	case Stderr:
		fi, err = os.Stderr.Stat()
	default:
		// Use Lstat here to not follow symlinks.
		// We will resolve symlinks target while this object's type is link.
		fi, err = os.Lstat(absPath)
	}

	return
}

func (s *Storage) getAbsPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(s.workDir, path)
}

func (s *Storage) formatError(op string, err error, path ...string) error {
	if err == nil {
		return nil
	}

	return &services.StorageError{
		Op:       op,
		Err:      formatError(err),
		Storager: s,
		Path:     path,
	}
}
