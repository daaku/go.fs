// Package emptyfs provides a File System that is always empty.
package emptyfs

import (
	"github.com/daaku/go.fs"
	"github.com/daaku/go.fs/fsutil"
)

type system struct {
	fixed error
}

var singleton = system{}

// Provides a File System that is always empty.
func New() fs.System {
	return singleton
}

// Provides a File System that always returns the provided error and considers
// it to to represent that the requested file does not exist.
func NewWithError(err error) fs.System {
	return system{fixed: err}
}

func (s system) Open(name string) (fs.File, error) {
	if s.fixed != nil {
		return nil, s.fixed
	}
	return nil, fsutil.NewErrNotFound(name)
}

func (s system) IsNotExist(err error) bool {
	if s.fixed != nil && s.fixed == err {
		return true
	}
	return fsutil.IsNotExist(err)
}
