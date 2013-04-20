// Package fs provides various File System related abstractions.
//
// It is often valuable to have a weak mapping between how we utilize files and
// directories in our code to the real file system.
//
// At first you just want the real file system: realfs. Then you write tests
// and just want to define the file system in test code: memfs. You get more
// developers on your project and need to reference files in relative terms and
// possibly limit access: limitfs. Your project becomes mature and you want
// easier deployent and start packaging your resources into a zip file: zipfs.
// Obviously you didn't just want a zip file, you want to just augment the
// compiled binary you already deploy, but you still want to be a `go get`
// compatible package during development: pkgfs. You want a tool to do the last
// step for you: pkgfszip.
//
// How to use it:
//
// Write your libraries using the interfaces in this package, and then use the
// implementations from any of the above as your use case evolves. Most likely
// you want to use memfs in your tests, pkgfs in application code and
// pkgfszip to augment your binary before deployment and you'll get all the
// benefits of this family of packages.
//
// A note about read & write:
//
// Not all file systems are created equal, but it helps to treat them as such.
// For this reason the general abstraction provides includes read as well as
// write APIs. For file systems like realfs, memfs & limitfs this is great
// since those file systems do in fact provide write APIs. On the other hand
// zipfs is read-only. For such scenarios the implementation just returns
// errors when you try to use the write APIs. In practice this doesn't mean
// much and you can mostly just ignore the write APIs if you live in a read
// only world and want it's advantages or use the write APIs and not use
// abstractions like zipfs or pkgfs which don't make much sense with respect
// to writes.
package fs

import (
	"os"
)

// A File implements access to a single file or directory.
type File interface {
	// Close closes the File, rendering it unusable for I/O.
	Close() error

	// Chmod changes the mode of the file to mode.
	Chmod(mode os.FileMode) error

	// Chown changes the numeric uid and gid of the named file.
	Chown(uid, gid int) error

	// Get the owner UID.
	OwnerGID() (int, error)

	// Get the owner GID.
	OwnerUID() (int, error)

	// Read reads up to len(b) bytes from the File. It returns the number of bytes
	// read and an error, if any. EOF is signaled by a zero count with err set to
	// io.EOF.
	Read(b []byte) (n int, err error)

	// ReadAt reads len(b) bytes from the File starting at byte offset off. It
	// returns the number of bytes read and the error, if any. ReadAt always
	// returns a non-nil error when n < len(b). At end of file, that error is
	// io.EOF.
	ReadAt(b []byte, off int64) (n int, err error)

	// Returns the FileInfos of the files in the directory.
	Readdir(n int) (infos []os.FileInfo, err error)

	// Returns names of files in the directory.
	Readdirnames(n int) (names []string, err error)

	// Seek sets the offset for the next Read or Write on file to offset,
	// interpreted according to whence: 0 means relative to the origin of the
	// file, 1 means relative to the current offset, and 2 means relative to the
	// end. It returns the new offset and an error, if any.
	Seek(offset int64, whence int) (ret int64, err error)

	// Stat returns the FileInfo structure describing this File.
	Stat() (fi os.FileInfo, err error)

	// Sync the file.
	Sync() (err error)

	// Truncate changes the size of the file. It does not change the I/O offset.
	Truncate(size int64) error

	// Write writes len(b) bytes to the File. It returns the number of bytes
	// written and an error, if any.
	Write(b []byte) (ret int, err error)

	// WriteAt writes len(b) bytes to the File starting at byte offset off. It
	// returns the number of bytes written and an error, if any. WriteAt returns a
	// non-nil error when n != len(b).
	WriteAt(b []byte, off int64) (ret int, err error)

	// WriteString is like Write, but writes the contents of string s rather than
	// an array of bytes.
	WriteString(s string) (ret int, err error)
}

// A System implements access to a collection of named files.
type System interface {
	// Open a named file for reading.
	Open(name string) (File, error)

	// IsNotExist returns whether the error is known to report that a file does
	// not exist.
	IsNotExist(err error) bool
}
