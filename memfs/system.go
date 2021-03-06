package memfs

import (
	"os"
	"path/filepath"
	"time"

	"github.com/daaku/go.fs"
	"github.com/daaku/go.fs/emptyfs"
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

// Creates a fs.System backed by the given map. It expects only Files and will
// generate Directory entries automatically.
func NewWithFiles(files map[string]fs.File) fs.System {
	s := make(map[string]fs.File)
	var add func(fullpath string, file fs.File) error
	add = func(fullpath string, file fs.File) error {
		fi, err := file.Stat()
		if err != nil {
			return err
		}
		s[fullpath] = file
		parent := filepath.Dir(fullpath)
		if parent == "." {
			return nil
		}
		if parentdir := s[parent]; parentdir != nil {
			pf := parentdir.(*File)
			if err := pf.AddDirInfo(fi); err != nil {
				return err
			}
		} else {
			parentdir := NewDir(
				filepath.Base(parent), os.FileMode(755), time.Now(), []os.FileInfo{fi})
			if err := add(parent, parentdir); err != nil {
				return err
			}
		}
		return nil
	}

	for fullpath, file := range files {
		if err := add(fullpath, file); err != nil {
			return emptyfs.NewWithError(err)
		}
	}
	return NewSystem(s)
}
