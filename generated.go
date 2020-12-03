// Code generated by go generate via internal/cmd/service; DO NOT EDIT.
package fs

import (
	"context"
	"io"

	"github.com/aos-dev/go-storage/v2/pkg/credential"
	"github.com/aos-dev/go-storage/v2/pkg/endpoint"
	"github.com/aos-dev/go-storage/v2/pkg/httpclient"
	"github.com/aos-dev/go-storage/v2/services"
	. "github.com/aos-dev/go-storage/v2/types"
)

var _ credential.Provider
var _ endpoint.Provider
var _ Storager
var _ services.ServiceError
var _ httpclient.Options

// Type is the type for fs
const Type = "fs"

// Service available pairs.
const ()

// pairStorageNew is the parsed struct
type pairStorageNew struct {
	pairs []Pair

	// Required pairs
	// Optional pairs
	HasPairPolicy bool
	PairPolicy    PairPolicy
	HasWorkDir    bool
	WorkDir       string
	// Generated pairs
}

// parsePairStorageNew will parse Pair slice into *pairStorageNew
func parsePairStorageNew(opts []Pair) (*pairStorageNew, error) {
	result := &pairStorageNew{
		pairs: opts,
	}

	for _, v := range opts {
		switch v.Key {
		// Required pairs
		// Optional pairs
		case "pair_policy":
			result.HasPairPolicy = true
			result.PairPolicy = v.Value.(PairPolicy)
		case "work_dir":
			result.HasWorkDir = true
			result.WorkDir = v.Value.(string)
			// Generated pairs
		}
	}

	return result, nil
}

// pairStorageCopy is the parsed struct
type pairStorageCopy struct {
	pairs []Pair

	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageCopy will parse Pair slice into *pairStorageCopy
func (s *Storage) parsePairStorageCopy(opts []Pair) (*pairStorageCopy, error) {
	result := &pairStorageCopy{
		pairs: opts,
	}

	for _, v := range opts {
		switch v.Key {
		// Required pairs
		// Optional pairs
		// Generated pairs
		default:

			if s.pairPolicy.All || s.pairPolicy.Copy {
				return nil, services.NewPairUnsupportedError(v)
			}

		}
	}

	return result, nil
}

// pairStorageDelete is the parsed struct
type pairStorageDelete struct {
	pairs []Pair

	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageDelete will parse Pair slice into *pairStorageDelete
func (s *Storage) parsePairStorageDelete(opts []Pair) (*pairStorageDelete, error) {
	result := &pairStorageDelete{
		pairs: opts,
	}

	for _, v := range opts {
		switch v.Key {
		// Required pairs
		// Optional pairs
		// Generated pairs
		default:

			if s.pairPolicy.All || s.pairPolicy.Delete {
				return nil, services.NewPairUnsupportedError(v)
			}

		}
	}

	return result, nil
}

// pairStorageFetch is the parsed struct
type pairStorageFetch struct {
	pairs []Pair

	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageFetch will parse Pair slice into *pairStorageFetch
func (s *Storage) parsePairStorageFetch(opts []Pair) (*pairStorageFetch, error) {
	result := &pairStorageFetch{
		pairs: opts,
	}

	for _, v := range opts {
		switch v.Key {
		// Required pairs
		// Optional pairs
		// Generated pairs
		default:

			if s.pairPolicy.All || s.pairPolicy.Fetch {
				return nil, services.NewPairUnsupportedError(v)
			}

		}
	}

	return result, nil
}

// pairStorageListDir is the parsed struct
type pairStorageListDir struct {
	pairs []Pair

	// Required pairs
	// Optional pairs
	HasContinuationToken bool
	ContinuationToken    string
	// Generated pairs
}

// parsePairStorageListDir will parse Pair slice into *pairStorageListDir
func (s *Storage) parsePairStorageListDir(opts []Pair) (*pairStorageListDir, error) {
	result := &pairStorageListDir{
		pairs: opts,
	}

	for _, v := range opts {
		switch v.Key {
		// Required pairs
		// Optional pairs
		case "continuation_token":
			result.HasContinuationToken = true
			result.ContinuationToken = v.Value.(string)
		// Generated pairs
		default:

			if s.pairPolicy.All || s.pairPolicy.ListDir {
				return nil, services.NewPairUnsupportedError(v)
			}

		}
	}

	return result, nil
}

// pairStorageMetadata is the parsed struct
type pairStorageMetadata struct {
	pairs []Pair

	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageMetadata will parse Pair slice into *pairStorageMetadata
func (s *Storage) parsePairStorageMetadata(opts []Pair) (*pairStorageMetadata, error) {
	result := &pairStorageMetadata{
		pairs: opts,
	}

	for _, v := range opts {
		switch v.Key {
		// Required pairs
		// Optional pairs
		// Generated pairs
		default:

			if s.pairPolicy.All || s.pairPolicy.Metadata {
				return nil, services.NewPairUnsupportedError(v)
			}

		}
	}

	return result, nil
}

// pairStorageMove is the parsed struct
type pairStorageMove struct {
	pairs []Pair

	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageMove will parse Pair slice into *pairStorageMove
func (s *Storage) parsePairStorageMove(opts []Pair) (*pairStorageMove, error) {
	result := &pairStorageMove{
		pairs: opts,
	}

	for _, v := range opts {
		switch v.Key {
		// Required pairs
		// Optional pairs
		// Generated pairs
		default:

			if s.pairPolicy.All || s.pairPolicy.Move {
				return nil, services.NewPairUnsupportedError(v)
			}

		}
	}

	return result, nil
}

// pairStorageRead is the parsed struct
type pairStorageRead struct {
	pairs []Pair

	// Required pairs
	// Optional pairs
	HasOffset           bool
	Offset              int64
	HasReadCallbackFunc bool
	ReadCallbackFunc    func([]byte)
	HasSize             bool
	Size                int64
	// Generated pairs
}

// parsePairStorageRead will parse Pair slice into *pairStorageRead
func (s *Storage) parsePairStorageRead(opts []Pair) (*pairStorageRead, error) {
	result := &pairStorageRead{
		pairs: opts,
	}

	for _, v := range opts {
		switch v.Key {
		// Required pairs
		// Optional pairs
		case "offset":
			result.HasOffset = true
			result.Offset = v.Value.(int64)
		case "read_callback_func":
			result.HasReadCallbackFunc = true
			result.ReadCallbackFunc = v.Value.(func([]byte))
		case "size":
			result.HasSize = true
			result.Size = v.Value.(int64)
		// Generated pairs
		default:

			if s.pairPolicy.All || s.pairPolicy.Read {
				return nil, services.NewPairUnsupportedError(v)
			}

		}
	}

	return result, nil
}

// pairStorageStat is the parsed struct
type pairStorageStat struct {
	pairs []Pair

	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageStat will parse Pair slice into *pairStorageStat
func (s *Storage) parsePairStorageStat(opts []Pair) (*pairStorageStat, error) {
	result := &pairStorageStat{
		pairs: opts,
	}

	for _, v := range opts {
		switch v.Key {
		// Required pairs
		// Optional pairs
		// Generated pairs
		default:

			if s.pairPolicy.All || s.pairPolicy.Stat {
				return nil, services.NewPairUnsupportedError(v)
			}

		}
	}

	return result, nil
}

// pairStorageWrite is the parsed struct
type pairStorageWrite struct {
	pairs []Pair

	// Required pairs
	// Optional pairs
	HasOffset           bool
	Offset              int64
	HasReadCallbackFunc bool
	ReadCallbackFunc    func([]byte)
	HasSize             bool
	Size                int64
	// Generated pairs
	HasContentMd5   bool
	ContentMd5      string
	HasContentType  bool
	ContentType     string
	HasStorageClass bool
	StorageClass    string
}

// parsePairStorageWrite will parse Pair slice into *pairStorageWrite
func (s *Storage) parsePairStorageWrite(opts []Pair) (*pairStorageWrite, error) {
	result := &pairStorageWrite{
		pairs: opts,
	}

	for _, v := range opts {
		switch v.Key {
		// Required pairs
		// Optional pairs
		case "offset":
			result.HasOffset = true
			result.Offset = v.Value.(int64)
		case "read_callback_func":
			result.HasReadCallbackFunc = true
			result.ReadCallbackFunc = v.Value.(func([]byte))
		case "size":
			result.HasSize = true
			result.Size = v.Value.(int64)
		// Generated pairs
		case "content_md5":
			value, ok := s.convertWriteContentMd5(v.Value.(string))
			if ok {
				result.HasContentMd5 = true
				result.ContentMd5 = value
			} else {

				if s.pairPolicy.All || s.pairPolicy.Write || s.pairPolicy.WriteContentMd5 {
					return nil, services.NewPairUnsupportedError(v)
				}

			}
		case "content_type":
			value, ok := s.convertWriteContentType(v.Value.(string))
			if ok {
				result.HasContentType = true
				result.ContentType = value
			} else {

				if s.pairPolicy.All || s.pairPolicy.Write || s.pairPolicy.WriteContentType {
					return nil, services.NewPairUnsupportedError(v)
				}

			}
		case "storage_class":
			value, ok := s.convertWriteStorageClass(v.Value.(string))
			if ok {
				result.HasStorageClass = true
				result.StorageClass = value
			} else {

				if s.pairPolicy.All || s.pairPolicy.Write || s.pairPolicy.WriteStorageClass {
					return nil, services.NewPairUnsupportedError(v)
				}

			}
		default:

			if s.pairPolicy.All || s.pairPolicy.Write {
				return nil, services.NewPairUnsupportedError(v)
			}

		}
	}

	return result, nil
}

// Copy will copy an Object or multiple object in the service.
//
// This function will create a context by default.
func (s *Storage) Copy(src string, dst string, pairs ...Pair) (err error) {
	ctx := context.Background()
	return s.CopyWithContext(ctx, src, dst, pairs...)
}

// CopyWithContext will copy an Object or multiple object in the service.
func (s *Storage) CopyWithContext(ctx context.Context, src string, dst string, pairs ...Pair) (err error) {
	defer func() {
		err = s.formatError("copy", err, src, dst)
	}()
	var opt *pairStorageCopy
	opt, err = s.parsePairStorageCopy(pairs)
	if err != nil {
		return
	}

	return s.copy(ctx, src, dst, opt)
}

// Delete will delete an Object from service.
//
// This function will create a context by default.
func (s *Storage) Delete(path string, pairs ...Pair) (err error) {
	ctx := context.Background()
	return s.DeleteWithContext(ctx, path, pairs...)
}

// DeleteWithContext will delete an Object from service.
func (s *Storage) DeleteWithContext(ctx context.Context, path string, pairs ...Pair) (err error) {
	defer func() {
		err = s.formatError("delete", err, path)
	}()
	var opt *pairStorageDelete
	opt, err = s.parsePairStorageDelete(pairs)
	if err != nil {
		return
	}

	return s.delete(ctx, path, opt)
}

// Fetch will fetch from a given url to path.
//
// This function will create a context by default.
func (s *Storage) Fetch(path string, url string, pairs ...Pair) (err error) {
	ctx := context.Background()
	return s.FetchWithContext(ctx, path, url, pairs...)
}

// FetchWithContext will fetch from a given url to path.
func (s *Storage) FetchWithContext(ctx context.Context, path string, url string, pairs ...Pair) (err error) {
	defer func() {
		err = s.formatError("fetch", err, path, url)
	}()
	var opt *pairStorageFetch
	opt, err = s.parsePairStorageFetch(pairs)
	if err != nil {
		return
	}

	return s.fetch(ctx, path, url, opt)
}

// ListDir will return list a specific dir.
//
// This function will create a context by default.
func (s *Storage) ListDir(dir string, pairs ...Pair) (oi *ObjectIterator, err error) {
	ctx := context.Background()
	return s.ListDirWithContext(ctx, dir, pairs...)
}

// ListDirWithContext will return list a specific dir.
func (s *Storage) ListDirWithContext(ctx context.Context, dir string, pairs ...Pair) (oi *ObjectIterator, err error) {
	defer func() {
		err = s.formatError("list_dir", err, dir)
	}()
	var opt *pairStorageListDir
	opt, err = s.parsePairStorageListDir(pairs)
	if err != nil {
		return
	}

	return s.listDir(ctx, dir, opt)
}

// Metadata will return current storager's metadata.
//
// This function will create a context by default.
func (s *Storage) Metadata(pairs ...Pair) (meta *StorageMeta, err error) {
	ctx := context.Background()
	return s.MetadataWithContext(ctx, pairs...)
}

// MetadataWithContext will return current storager's metadata.
func (s *Storage) MetadataWithContext(ctx context.Context, pairs ...Pair) (meta *StorageMeta, err error) {
	defer func() {
		err = s.formatError("metadata", err)
	}()
	var opt *pairStorageMetadata
	opt, err = s.parsePairStorageMetadata(pairs)
	if err != nil {
		return
	}

	return s.metadata(ctx, opt)
}

// Move will move an object in the service.
//
// This function will create a context by default.
func (s *Storage) Move(src string, dst string, pairs ...Pair) (err error) {
	ctx := context.Background()
	return s.MoveWithContext(ctx, src, dst, pairs...)
}

// MoveWithContext will move an object in the service.
func (s *Storage) MoveWithContext(ctx context.Context, src string, dst string, pairs ...Pair) (err error) {
	defer func() {
		err = s.formatError("move", err, src, dst)
	}()
	var opt *pairStorageMove
	opt, err = s.parsePairStorageMove(pairs)
	if err != nil {
		return
	}

	return s.move(ctx, src, dst, opt)
}

// Read will read the file's data.
//
// This function will create a context by default.
func (s *Storage) Read(path string, w io.Writer, pairs ...Pair) (n int64, err error) {
	ctx := context.Background()
	return s.ReadWithContext(ctx, path, w, pairs...)
}

// ReadWithContext will read the file's data.
func (s *Storage) ReadWithContext(ctx context.Context, path string, w io.Writer, pairs ...Pair) (n int64, err error) {
	defer func() {
		err = s.formatError("read", err, path)
	}()
	var opt *pairStorageRead
	opt, err = s.parsePairStorageRead(pairs)
	if err != nil {
		return
	}

	return s.read(ctx, path, w, opt)
}

// Stat will stat a path to get info of an object.
//
// This function will create a context by default.
func (s *Storage) Stat(path string, pairs ...Pair) (o *Object, err error) {
	ctx := context.Background()
	return s.StatWithContext(ctx, path, pairs...)
}

// StatWithContext will stat a path to get info of an object.
func (s *Storage) StatWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error) {
	defer func() {
		err = s.formatError("stat", err, path)
	}()
	var opt *pairStorageStat
	opt, err = s.parsePairStorageStat(pairs)
	if err != nil {
		return
	}

	return s.stat(ctx, path, opt)
}

// Write will write data into a file.
//
// This function will create a context by default.
func (s *Storage) Write(path string, r io.Reader, pairs ...Pair) (n int64, err error) {
	ctx := context.Background()
	return s.WriteWithContext(ctx, path, r, pairs...)
}

// WriteWithContext will write data into a file.
func (s *Storage) WriteWithContext(ctx context.Context, path string, r io.Reader, pairs ...Pair) (n int64, err error) {
	defer func() {
		err = s.formatError("write", err, path)
	}()
	var opt *pairStorageWrite
	opt, err = s.parsePairStorageWrite(pairs)
	if err != nil {
		return
	}

	return s.write(ctx, path, r, opt)
}
