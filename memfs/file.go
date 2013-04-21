package memfs

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"
)

var (
	errNotDir             = errors.New("memfs: file is not a directory")
	errIsDir              = errors.New("memfs: file is a directory")
	errAlreadyClosed      = errors.New("memfs: file already closed")
	errTruncateOutOfRange = errors.New("memfs: file truncation out of range")
	errSeekOutOfRange     = errors.New("memfs: file seek out of range")
	errSeekWhenceInvalid  = errors.New("memfs: file seek whence invalid")
	errOffsetOutOfRange   = errors.New("memfs: file offset out of range")
)

// In-memory File representation.
type File struct {
	name     string
	fileInfo *MemFileInfo
	uid      int
	gid      int
	closed   bool
	isDir    bool
	off      int64         // dual purpose for buf & infos depending on isDir
	buf      []byte        // for files
	infos    []os.FileInfo // for directories
}

// Create a new File.
func NewFile(name string, mode os.FileMode, mtime time.Time, data []byte) *File {
	return &File{
		name: name,
		buf:  data,
		fileInfo: NewFileInfo(FileInfo{
			Name:    filepath.Base(name),
			Size:    int64(len(data)),
			Mode:    mode,
			ModTime: mtime,
		}),
	}
}

// Create a new Directory.
func NewDir(name string, mode os.FileMode, mtime time.Time, infos []os.FileInfo) *File {
	return &File{
		isDir: true,
		name:  name,
		infos: infos,
		fileInfo: NewFileInfo(FileInfo{
			Name:    filepath.Base(name),
			Mode:    mode,
			ModTime: mtime,
		}),
	}
}

// Chmod changes the mode of the file to mode.
func (f *File) Chmod(mode os.FileMode) error {
	f.fileInfo.SetMode(mode)
	return nil
}

// Chown changes the numeric uid and gid of the named file.
func (f *File) Chown(uid, gid int) error {
	f.uid = uid
	f.gid = gid
	return nil
}

// Get the owner UID.
func (f *File) OwnerUID() (int, error) {
	return f.uid, nil
}

// Get the owner GID.
func (f *File) OwnerGID() (int, error) {
	return f.gid, nil
}

// Close closes the File, rendering it unusable for I/O.
func (f *File) Close() error {
	f.closed = true
	return nil
}

// Check if the File has been closed.
func (f *File) IsClosed() bool {
	return f.closed
}

// Name returns the name of the file as presented to Open.
func (f *File) Name() string {
	return f.name
}

// Name returns the name of the file as presented to Open.
func (f *File) SetName(name string) {
	f.name = name
	f.fileInfo.SetName(filepath.Base(name))
}

// Read reads up to len(b) bytes from the File. It returns the number of bytes
// read and an error, if any. EOF is signaled by a zero count with err set to
// io.EOF.
func (f *File) Read(b []byte) (n int, err error) {
	if f.IsClosed() {
		return 0, errAlreadyClosed
	}

	if f.isDir {
		return 0, errIsDir
	}

	if f.off >= int64(len(f.buf)) {
		if len(b) == 0 {
			return
		}
		return 0, io.EOF
	}
	n = copy(b, f.buf[f.off:])
	f.off += int64(n)
	return
}

// ReadAt reads len(b) bytes from the File starting at byte offset off. It
// returns the number of bytes read and the error, if any. ReadAt always
// returns a non-nil error when n < len(b). At end of file, that error is
// io.EOF.
func (f *File) ReadAt(b []byte, off int64) (n int, err error) {
	if f.IsClosed() {
		return 0, errAlreadyClosed
	}

	if f.isDir {
		return 0, errIsDir
	}

	f.off = off
	return f.Read(b)
}

// Returns the FileInfos of the files in the directory.
func (f *File) Readdir(n int) (infos []os.FileInfo, err error) {
	if !f.isDir {
		return nil, errNotDir
	}

	if n <= 0 {
		f.off = 0
		return f.infos, nil
	}

	m := f.off + int64(n)
	l := int64(len(f.infos))
	if m >= l {
		err = io.EOF
		m = l
	}

	infos = f.infos[f.off:m]
	f.off = m
	return
}

// Returns names of files in the directory.
func (f *File) Readdirnames(n int) (names []string, err error) {
	if !f.isDir {
		return nil, errNotDir
	}

	fis, err := f.Readdir(n)
	names = make([]string, len(fis))
	for ix, fi := range fis {
		names[ix] = fi.Name()
	}
	return
}

// Reset infos for a directory.
func (f *File) SetDirInfos(infos []os.FileInfo) error {
	if !f.isDir {
		return errNotDir
	}

	f.infos = infos
	f.Reset()
	return nil
}

// Add a new info to the directory. Will also reset the internal offset.
func (f *File) AddDirInfo(info os.FileInfo) error {
	if !f.isDir {
		return errNotDir
	}

	f.infos = append(f.infos, info)
	f.Reset()
	return nil
}

// Reset offset for Read/Write/Readdir/Readdirnames.
func (f *File) Reset() {
	f.off = 0
}

// Seek sets the offset for the next Read or Write on file to offset,
// interpreted according to whence: 0 means relative to the origin of the file,
// 1 means relative to the current offset, and 2 means relative to the end. It
// returns the new offset and an error, if any.
func (f *File) Seek(offset int64, whence int) (ret int64, err error) {
	if f.IsClosed() {
		return 0, errAlreadyClosed
	}

	if f.isDir {
		return 0, errIsDir
	}

	l := int64(len(f.buf))
	switch whence {
	case os.SEEK_SET:
		ret = offset
	case os.SEEK_CUR:
		ret = f.off + offset
	case os.SEEK_END:
		ret = l - offset
	default:
		return f.off, errSeekWhenceInvalid
	}

	if ret > l || ret < 0 {
		return f.off, errSeekOutOfRange
	}
	f.off = ret
	return ret, nil
}

// Stat returns the FileInfo structure describing this File.
func (f *File) Stat() (fi os.FileInfo, err error) {
	if f.IsClosed() {
		return nil, errAlreadyClosed
	}

	return f.fileInfo, nil
}

// For in memory files Sync does nothing.
func (f *File) Sync() (err error) {
	if f.IsClosed() {
		return errAlreadyClosed
	}

	return nil
}

// Truncate changes the size of the file. It does not change the I/O offset.
func (f *File) Truncate(size int64) error {
	if f.IsClosed() {
		return errAlreadyClosed
	}

	if f.isDir {
		return errIsDir
	}

	if size < 0 || size > int64(len(f.buf)) {
		return errTruncateOutOfRange
	}
	f.buf = f.buf[0:size]
	f.updateFileInfoSize()
	return nil
}

// Write writes len(b) bytes to the File. It returns the number of bytes
// written and an error, if any.
func (f *File) Write(b []byte) (ret int, err error) {
	if f.IsClosed() {
		return 0, errAlreadyClosed
	}

	if f.isDir {
		return 0, errIsDir
	}

	f.grow(len(b))
	ret = copy(f.buf[f.off:], b)
	f.off += int64(ret)
	f.updateFileInfoSize()
	return ret, nil
}

// WriteAt writes len(b) bytes to the File starting at byte offset off. It
// returns the number of bytes written and an error, if any. WriteAt returns a
// non-nil error when n != len(b).
func (f *File) WriteAt(b []byte, off int64) (ret int, err error) {
	if f.IsClosed() {
		return 0, errAlreadyClosed
	}

	if f.isDir {
		return 0, errIsDir
	}

	if off > int64(len(f.buf)) {
		return 0, errOffsetOutOfRange
	}
	f.grow(len(b))
	ret = copy(f.buf[off:], b)
	f.updateFileInfoSize()
	return ret, nil
}

// WriteString is like Write, but writes the contents of string s rather than
// an array of bytes.
func (f *File) WriteString(s string) (ret int, err error) {
	if f.IsClosed() {
		return 0, errAlreadyClosed
	}

	if f.isDir {
		return 0, errIsDir
	}

	f.grow(len(s))
	ret = copy(f.buf[f.off:], s)
	f.off += int64(ret)
	f.updateFileInfoSize()
	return ret, nil
}

// Grows the buffer to guarantee space for n more bytes. Note, this modifies
// the length of the buffer in addition to the capacity to allow for a simple
// copy operation to follow.
func (f *File) grow(n int) {
	l := len(f.buf)
	c := cap(f.buf)
	if l+n > c {
		buf := make([]byte, (2*c)+n)
		copy(buf, f.buf)
		f.buf = buf[:l+n]
	} else {
		f.buf = f.buf[:l+n]
	}
}

// Updates the size in the underlying FileInfo.
func (f *File) updateFileInfoSize() {
	f.fileInfo.SetSize(int64(len(f.buf)))
}
