// Package zipfs provides a File System interface backed by a zip file.
package zipfs

import (
	"archive/zip"
	"errors"
	"io"
	"os"

	"github.com/daaku/go.fs"
	"github.com/daaku/go.fs/fsutil"
	"github.com/daaku/go.zipexe"
)

type file struct {
	io.ReadCloser
	*zip.File
}

func (f *file) Stat() (os.FileInfo, error) {
	return f.FileInfo(), nil
}

func (f *file) Chmod(mode os.FileMode) error {
	return errors.New("zipfs: Chmod not supported on file")
}

func (f *file) Chown(uid, gid int) error {
	return errors.New("zipfs: Chown not supported on file")
}

func (f *file) OwnerGID() (int, error) {
	return 0, errors.New("zipfs: OwnerGID not supported on file")
}

func (f *file) OwnerUID() (int, error) {
	return 0, errors.New("zipfs: OwnerUID not supported on file")
}

func (f *file) ReadAt(b []byte, off int64) (n int, err error) {
	return 0, errors.New("zipfs: ReadAt not supported on file")
}

func (f *file) Readdir(count int) ([]os.FileInfo, error) {
	return nil, errors.New("zipfs: Readdir not supported on file")
}

func (f *file) Readdirnames(n int) (names []string, err error) {
	return nil, errors.New("zipfs: Readdirnames not supported on file")
}

func (f *file) Seek(offset int64, whence int) (ret int64, err error) {
	return 0, errors.New("zipfs: Seek not supported on file")
}

func (f *file) Sync() (err error) {
	return nil
}

func (f *file) Truncate(size int64) error {
	return errors.New("zipfs: Truncate not supported on file")
}

func (f *file) Write(b []byte) (ret int, err error) {
	return 0, errors.New("zipfs: Write not supported on file")
}

func (f *file) WriteAt(b []byte, off int64) (ret int, err error) {
	return 0, errors.New("zipfs: WriteAt not supported on file")
}

func (f *file) WriteString(s string) (ret int, err error) {
	return 0, errors.New("zipfs: WriteString not supported on file")
}

type system struct {
	zipReader *zip.Reader
}

func (s *system) Open(name string) (fs.File, error) {
	for _, f := range s.zipReader.File {
		if f.Name == name {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			return &file{
				ReadCloser: rc,
				File:       f,
			}, nil
		}
	}
	return nil, fsutil.NewErrNotFound(name)
}

func (s *system) IsNotExist(err error) bool {
	return fsutil.IsNotExist(err)
}

// Open a file system using the given zip.Reader.
func New(zr *zip.Reader) fs.System {
	return &system{zipReader: zr}
}

// Opens the named zip file as a fs.System.
func Open(name string) (fs.System, error) {
	zr, err := zipexe.Open(name)
	if err != nil {
		return nil, err
	}
	return New(zr), nil
}
