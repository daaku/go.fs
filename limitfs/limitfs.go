// Package limitfs allows for providing a limited view on an existing
// fs.System source.
package limitfs

import (
	"github.com/daaku/go.fs"
	"os"
)

// Defines a Config that selects files to be made available via a File System.
type Config struct {
	Root      string // used as the root of the File System
	Recursive bool   // control access to nested directories
	Glob      string // limit by a glob pattern
}

type system struct {
	fs.File
	System    fs.System
	glob      string
	recusrive bool
}

func (s system) Readdir(count int) (fis []os.FileInfo, err error) {
	if count <= 0 {
		raw, err := s.File.Readdir(count)
		fis, errF := s.filter(raw)
		if err == nil && errF != nil {
			err = errF
		}
		return fis, err
	}

	pending := count
	for pending > 0 {
		raw, err := s.File.Readdir(pending)
		if len(raw) == 0 {
			return fis, err
		}
		filtered, errF := s.filter(raw)
		fis = append(fis, filtered...)
		if err == nil && errF != nil {
			err = errF
		}
		if err != nil {
			return fis, err
		}
		pending = count - len(fis)
	}
	return
}

func (s system) filter(given []os.FileInfo) (final []os.FileInfo, err error) {
	// FIXME
	for _, fi := range given {
		if fi.IsDir() {
			if !s.recusrive {
				continue
			}
			//cfs, err := 1, nil
		}
		final = append(final, fi)
	}
	return final, nil
}

// Create a wrapped fs.System that limits access based on the provided Config.
func New(c Config, system fs.System) fs.System {
	return nil
}
