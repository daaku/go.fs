// Package emptyfs provides a File System that is always empty.
package emptyfs

import (
	"github.com/daaku/go.fs"
	"github.com/daaku/go.fs/fsutil"
)

type system struct{}

var singleton = system{}

// Provides a File System that is always empty.
func New() fs.System {
	return singleton
}

func (s system) Open(name string) (fs.File, error) {
	return nil, fsutil.NewErrNotFound(name)
}

func (s system) IsNotExist(err error) bool {
	return fsutil.IsNotExist(err)
}
