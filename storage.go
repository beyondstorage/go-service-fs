package fs

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/qingstor/go-mime"

	"github.com/aos-dev/go-storage/v2/pkg/iowrap"
	typ "github.com/aos-dev/go-storage/v2/types"
)

func (s *Storage) delete(ctx context.Context, path string, opt *pairStorageDelete) (err error) {
	rp := s.getAbsPath(path)

	err = os.Remove(rp)
	if err != nil {
		return err
	}
	return nil
}

type listDirInput struct {
	rp  string
	dir string

	enableLinkFollow bool

	started           bool
	continuationToken string

	f   *os.File
	buf []byte
}

func (input *listDirInput) ContinuationToken() string {
	return input.continuationToken
}

func (s *Storage) listDir(ctx context.Context, dir string, opt *pairStorageListDir) (oi *typ.ObjectIterator, err error) {
	input := listDirInput{
		// Always keep service original name as rp.
		rp: s.getAbsPath(dir),
		// Then convert the dir to slash separator.
		dir:              filepath.ToSlash(dir),
		enableLinkFollow: opt.EnableLinkFollow,

		// if HasContinuationToken, we should start after we scanned this token.
		// else, we can start directly.
		started:           !opt.HasContinuationToken,
		continuationToken: opt.ContinuationToken,

		buf: make([]byte, 8192),
	}

	return typ.NewObjectIterator(ctx, s.listDirNext, &input), nil
}

func (s *Storage) metadata(ctx context.Context, opt *pairStorageMetadata) (meta typ.StorageMeta, err error) {
	meta = typ.NewStorageMeta()
	meta.WorkDir = s.workDir
	return meta, nil
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt *pairStorageRead) (n int64, err error) {
	var rc io.ReadCloser

	rp := s.getAbsPath(path)

	f, err := s.openFile(rp)
	if err != nil {
		return
	}
	if opt.HasOffset {
		_, err = f.Seek(opt.Offset, 0)
		if err != nil {
			return n, err
		}
	}

	rc = f

	if opt.HasSize {
		rc = iowrap.LimitReadCloser(rc, opt.Size)
	}
	if opt.HasReadCallbackFunc {
		rc = iowrap.CallbackReadCloser(rc, opt.ReadCallbackFunc)
	}

	return io.Copy(w, rc)
}

func (s *Storage) stat(ctx context.Context, path string, opt *pairStorageStat) (o *typ.Object, err error) {
	rp := s.getAbsPath(path)

	fi, err := s.statFile(rp)
	if err != nil {
		return nil, err
	}

	o = s.newObject(true)
	o.ID = rp
	o.Name = path

	if fi.IsDir() {
		o.Type = typ.ObjectTypeDir
		return
	}

	if fi.Mode().IsRegular() {
		o.SetSize(fi.Size())
		o.SetUpdatedAt(fi.ModTime())

		if v := mime.DetectFilePath(path); v != "" {
			o.SetContentType(v)
		}

		o.Type = typ.ObjectTypeFile
		return
	}
	if fi.Mode()&StreamModeType != 0 {
		o.Type = typ.ObjectTypeStream
		return
	}

	o.Type = typ.ObjectTypeInvalid
	return o, nil
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, opt *pairStorageWrite) (n int64, err error) {
	var f io.WriteCloser

	rp := s.getAbsPath(path)

	f, err = s.createFile(rp)
	if err != nil {
		return
	}

	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReader(r, opt.ReadCallbackFunc)
	}

	if opt.HasSize {
		return io.CopyN(f, r, opt.Size)
	}
	return io.CopyBuffer(f, r, make([]byte, 1024*1024))
}

func (s *Storage) copy(ctx context.Context, src string, dst string, opt *pairStorageCopy) (err error) {
	rs := s.getAbsPath(src)
	rd := s.getAbsPath(dst)

	srcFile, err := s.openFile(rs)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := s.createFile(rd)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.CopyBuffer(dstFile, srcFile, make([]byte, 1024*1024))
	if err != nil {
		return err
	}
	return
}
func (s *Storage) move(ctx context.Context, src string, dst string, opt *pairStorageMove) (err error) {
	rs := s.getAbsPath(src)
	rd := s.getAbsPath(dst)

	// Create dir for dst path.
	// Create dir before create file
	err = os.MkdirAll(filepath.Dir(rd), 0755)
	if err != nil {
		return err
	}

	err = os.Rename(rs, rd)
	if err != nil {
		return err
	}
	return
}

func checkLink(v os.FileInfo, dir string) (os.FileInfo, error) {
	// if v is not link, return directly
	if v.Mode()&os.ModeSymlink == 0 {
		return v, nil
	}

	// otherwise, follow the link to get the target
	tarPath, err := filepath.EvalSymlinks(filepath.Join(dir, v.Name()))
	if err != nil {
		return nil, err
	}
	return os.Stat(tarPath)
}
