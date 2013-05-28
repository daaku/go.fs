// Package limitfs provides a view of another File System limiting access per
// your configuration.
package limitfs

import (
	"os"
	"path"
	"strings"

	"github.com/daaku/go.fs"
	"github.com/daaku/go.fs/fsutil"
)

// Defines a Config that selects files to be made available via a File System.
type Config struct {
	Root      string // used as the root of the File System
	Recursive bool   // control access to nested directories
	Glob      string // limit by a glob pattern
}

type system struct {
	Config Config
	System fs.System
}

func (s system) Open(name string) (fs.File, error) {
	cleaned, err := fsutil.Clean(name)
	if err != nil {
		return nil, err
	}

	if !s.Config.Recursive && strings.ContainsRune(cleaned, '/') {
		return nil, fsutil.NewErrLimitedNotFound(name)
	}

	final := path.Join(s.Config.Root, cleaned)
	f, err := s.System.Open(final)
	if err != nil {
		return nil, err
	}

	if s.Config.Glob != "" {
		fi, err := f.Stat()
		if err != nil {
			return nil, err
		}
		if fi.IsDir() {
			return dir{
				File: f,
				Path: final,
				Glob: s.Config.Glob,
			}, nil
		}
		match, err := path.Match(s.Config.Glob, final)
		if err != nil {
			return nil, err
		}
		if !match {
			return nil, fsutil.NewErrLimitedNotFound(name)
		}
	}

	return f, nil
}

func (s system) IsNotExist(err error) bool {
	return fsutil.IsNotExist(err)
}

type dir struct {
	fs.File
	Path string
	Glob string
}

func (d dir) Readdir(count int) (fis []os.FileInfo, err error) {
	if count <= 0 {
		raw, err := d.File.Readdir(count)
		fis, errF := d.filter(raw)
		if err == nil && errF != nil {
			err = errF
		}
		return fis, err
	}

	pending := count
	for pending > 0 {
		raw, err := d.File.Readdir(pending)
		filtered, errF := d.filter(raw)
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

func (d dir) filter(given []os.FileInfo) (final []os.FileInfo, err error) {
	for _, fi := range given {
		p := path.Join(d.Path, fi.Name())
		match, err := path.Match(d.Glob, p)
		if err != nil {
			return nil, err
		}
		if match {
			final = append(final, fi)
		}
	}
	return final, nil
}

func (d dir) Readdirnames(count int) (names []string, err error) {
	fis, err := d.Readdir(count)
	for _, fi := range fis {
		names = append(names, fi.Name())
	}
	return names, err
}

// Create a wrapped fs.System that limits access based on the provided Config.
func New(c Config, s fs.System) fs.System {
	return system{Config: c, System: s}
}
