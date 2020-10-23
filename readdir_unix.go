// +build aix darwin dragonfly freebsd js,wasm linux netbsd openbsd solaris

package fs

import (
	"context"
	"path"
	"path/filepath"

	"golang.org/x/sys/unix"

	typ "github.com/aos-dev/go-storage/v2/types"
)

// Available value for Dirent Type
//
// Copied from linux kernel <dirent.h>
// #define	DT_UNKNOWN	 0
// #define	DT_FIFO		 1
// #define	DT_CHR		 2
// #define	DT_DIR		 4
// #define	DT_BLK		 6
// #define	DT_REG		 8
// #define	DT_LNK		10
// #define	DT_SOCK		12
// #define	DT_WHT		14
const (
	// The file type is unknown.
	DirentTypeUnknown = 0
	// This is a named pipe (FIFO).
	DirentTypeFIFO = 1
	// This is a character device.
	DirentTypeCharDevice = 2
	// This is a directory.
	DirentTypeDirectory = 4
	// This is a block device.
	DirentTypeBlockDevice = 6
	// This is a regular file.
	DirentTypeRegular = 8
	// This is a symbolic link.
	DirentTypeLink = 10
	// This is a UNIX domain socket.
	DirentTypeSocket = 12
	// WhiteOut from BSD, don't know what's it mean.
	DirentTypeWhiteOut = 14
)

func getFiles(fd int, buf []byte) (files []file, err error) {
	n, err := unix.ReadDirent(fd, buf)
	if err != nil {
		return nil, err
	}
	if n <= 0 {
		return nil, nil
	}

	for len(buf) > 0 {
		// Get and check reclen
		reclen, ok := direntReclen(buf)
		if !ok || reclen > uint64(len(buf)) {
			return
		}
		// current dirent
		rec := buf[:reclen]
		// remaining dirents
		buf = buf[reclen:]

		// Get and check inode
		ino, ok := direntIno(rec)
		if !ok {
			break
		}
		if ino == 0 { // File absent in directory.
			continue
		}

		// Get and check type
		ty, ok := direntType(rec)
		if !ok {
			continue
		}

		// Get and check name
		name := rec[direntOffsetName:reclen]
		for i, c := range name {
			if c == 0 {
				name = name[:i]
				break
			}
		}
		// Check for useless names before allocating a string.
		if string(name) == "." || string(name) == ".." {
			continue
		}

		files = append(files, file{
			name: string(name),
			ty:   ty,
		})
	}

	return files, nil
}

type file struct {
	name string
	ty   uint8
}

func (s *Storage) listDirNext(ctx context.Context, page *typ.ObjectPage) (err error) {
	input := page.Status.(*listDirInput)

	// Open dir before we read it.
	if input.f == nil {
		input.f, err = s.osOpen(input.rp)
		if err != nil {
			return
		}
	}

	files, err := getFiles(int(input.f.Fd()), input.buf)
	if err != nil {
		return
	}

	// Whole dir has been read, return IterateDone to mark this iteration is done
	if len(files) == 0 {
		_ = input.f.Close()
		input.f = nil
		return typ.IterateDone
	}

	for _, v := range files {
		o := s.newObject(false)
		// Always keep service original name as ID.
		o.ID = filepath.Join(input.rp, v.name)
		// Object's name should always be separated by slash (/)
		o.Name = path.Join(input.dir, v.name)

		switch v.ty {
		case DirentTypeDirectory:
			o.Type = typ.ObjectTypeDir
		case DirentTypeRegular:
			o.Type = typ.ObjectTypeFile
		}

		page.Data = append(page.Data, o)
	}

	return
}
