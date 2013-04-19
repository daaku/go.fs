// Package realfs provides access to the real File System.
package realfs

import (
	"errors"
	"github.com/daaku/go.fs"
	"os"
)

type system struct{}

var singleton = system{}

// Provides access to the real unmodified file system.
func New() fs.System {
	return singleton
}

func (s system) Open(name string) (fs.File, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return file{f}, nil
}

func (s system) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

type file struct {
	*os.File
}

func (f file) OwnerGID() (int, error) {
	return 0, errors.New("realfs: TODO: OwnerGID not implemented")
}

func (f file) OwnerUID() (int, error) {
	return 0, errors.New("realfs: TODO: OwnerUID not implemented")
}
