package fs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/beyondstorage/go-storage/v4/pkg/httpclient"
	"github.com/beyondstorage/go-storage/v4/services"
	typ "github.com/beyondstorage/go-storage/v4/types"
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
	features     StorageFeatures

	typ.UnimplementedStorager
	typ.UnimplementedCopier
	typ.UnimplementedMover
	typ.UnimplementedFetcher
	typ.UnimplementedAppender
	typ.UnimplementedDirer
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
			err = services.InitError{Op: "new_storager", Type: Type, Err: formatError(err), Pairs: pairs}
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
	if opt.HasStorageFeatures {
		store.features = opt.StorageFeatures
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
	if _, ok := err.(services.InternalError); ok {
		return err
	}

	log.Printf("got error: %#+v", err)

	// Handle path & link errors.
	switch ie := err.(type) {
	case *os.PathError:
		switch ie.Err {
		case syscall.EISDIR:
			return fmt.Errorf("%w: %v", services.ErrObjectModeInvalid, err)
		}
	case *os.LinkError:
		log.Printf("got link error: %s", ie.Unwrap())
		switch ie.Err {
		// Golang will return syscall.EEXIST when move dst is dir on unix
		//
		// ref: https://golang.org/src/os/file_unix.go#38
		//
		// FIXME: maybe we need to move this part into utils_unix.go instead?
		case syscall.EEXIST:
			return fmt.Errorf("%w: %v", services.ErrObjectModeInvalid, err)
		// Golang will return syscall.EAFNOSUPPORT when move dst is dir on windows
		//
		// FIXME: maybe we need to move this part into utils_windows.go instead?
		case syscall.EAFNOSUPPORT:
			return fmt.Errorf("%w: %v", services.ErrObjectModeInvalid, err)
		}
	}

	// Handle error returned by os package.
	switch {
	case errors.Is(err, os.ErrNotExist):
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	case errors.Is(err, os.ErrPermission):
		return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
	default:
		return fmt.Errorf("%w: %v", services.ErrUnexpected, err)
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

	return services.StorageError{
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
