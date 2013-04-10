package memfs

import (
	"os"
	"time"
)

// Literal definition of a FileInfo that can be converted to an os.FileInfo.
// This provides a convinent way to define them.
type FileInfo struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	Sys     interface{}
}

// Since the interface and field names conflict, we copy over the data to this
// inernal struct that satisfies the public os.FileInfo interface.
type MemFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	sys     interface{}
}

// Create a new in-memory FileInfo based on the configured FileInfo.
func NewFileInfo(fi FileInfo) *MemFileInfo {
	return &MemFileInfo{
		name:    fi.Name,
		size:    fi.Size,
		mode:    fi.Mode,
		modTime: fi.ModTime,
		sys:     fi.Sys,
	}
}

// Create a new in-memory FileInfo based on a os.FileInfo.
func NewFileInfoFromExisting(fi os.FileInfo) *MemFileInfo {
	return &MemFileInfo{
		name:    fi.Name(),
		size:    fi.Size(),
		mode:    fi.Mode(),
		modTime: fi.ModTime(),
		sys:     fi.Sys(),
	}
}

// Base name for file.
func (fi *MemFileInfo) Name() string {
	return fi.name
}

// Set base name for file.
func (fi *MemFileInfo) SetName(name string) {
	fi.name = name
}

// Length in bytes for file.
func (fi *MemFileInfo) Size() int64 {
	return fi.size
}

// Set length in bytes for file.
func (fi *MemFileInfo) SetSize(size int64) {
	fi.size = size
}

// File mode bits for file.
func (fi *MemFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Set file mode bits for file.
func (fi *MemFileInfo) SetMode(mode os.FileMode) {
	fi.mode = mode
}

// Modification time for file.
func (fi *MemFileInfo) ModTime() time.Time {
	return fi.modTime
}

// Set modification time for file.
func (fi *MemFileInfo) SetModTime(t time.Time) {
	fi.modTime = t
}

// Abbreviation for Mode().IsDir().
func (fi *MemFileInfo) IsDir() bool {
	return fi.mode.IsDir()
}

// System specific data.
func (fi *MemFileInfo) Sys() interface{} {
	return fi.sys
}

// Set system specific data.
func (fi *MemFileInfo) SetSys(sys interface{}) {
	fi.sys = sys
}
