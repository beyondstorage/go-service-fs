package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aos-dev/go-storage/v3/pkg/httpclient"
	"github.com/aos-dev/go-storage/v3/services"
	typ "github.com/aos-dev/go-storage/v3/types"
)

const (
	// Std{in/out/err} support
	Stdin  = "/dev/stdin"
	Stdout = "/dev/stdout"
	Stderr = "/dev/stderr"

	PathSeparator = string(filepath.Separator)
)

// Storage is the fs client.
type Storage struct {
	// options for this storager.
	workDir string // workDir dir for all operation.

	defaultPairs DefaultStoragePairs
	pairPolicy   typ.PairPolicy

	typ.UnimplementedStorager
	typ.UnimplementedCopier
	typ.UnimplementedMover
	typ.UnimplementedFetcher
	typ.UnimplementedAppender
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

	if opt.HasDefaultStoragePairs {
		store.defaultPairs = opt.DefaultStoragePairs
	}
	if opt.HasPairPolicy {
		store.pairPolicy = opt.PairPolicy
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

func (s *Storage) openFile(absPath string, mode int) (f *os.File, needClose bool, err error) {
	switch absPath {
	case Stdin:
		f = os.Stdin
	case Stdout:
		f = os.Stdout
	case Stderr:
		f = os.Stderr
	default:
		needClose = true
		f, err = os.OpenFile(absPath, mode, 0664)
	}

	return
}

func (s *Storage) createFile(absPath string) (f *os.File, needClose bool, err error) {
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
			return nil, false, err
		}

		needClose = true
		f, err = os.Create(absPath)
	}

	return
}

func (s *Storage) createDir(absPath string) (err error) {
	defer func() {
		err = s.formatError("create_dir", err, absPath)
	}()

	return os.MkdirAll(absPath, 0755)
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
	absPath := filepath.Join(s.workDir, path)

	// Join will clean the trailing "/", we need to append it back.
	if strings.HasSuffix(path, PathSeparator) {
		absPath += PathSeparator
	}
	return absPath
}

func (s *Storage) isDirPath(path string) bool {
	return strings.HasSuffix(path, PathSeparator)
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

func (s *Storage) convertWriteContentMd5(v string) (string, bool) {
	return "", true
}

func (s *Storage) convertWriteContentType(v string) (string, bool) {
	return "", true
}

func (s *Storage) convertWriteStorageClass(v string) (string, bool) {
	return "", true
}

func convertNewHTTPClientOptions(_ *httpclient.Options) (*httpclient.Options, bool) {
	return nil, false
}
