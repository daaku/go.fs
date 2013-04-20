package memfs

import (
	"github.com/daaku/go.fs"
	"github.com/daaku/go.fs/fsutil"
)

type system map[string]fs.File

func (s system) Open(name string) (fs.File, error) {
	if f := s[name]; f != nil {
		return f, nil
	}
	return nil, fsutil.NewErrNotFound(name)
}

func (s system) IsNotExist(err error) bool {
	return fsutil.IsNotExist(err)
}

// Creates a fs.System backed by the given map. It expects directories to also
// have provided entries as necessary and won't create them.
func NewSystem(files map[string]fs.File) fs.System {
	return system(files)
}
