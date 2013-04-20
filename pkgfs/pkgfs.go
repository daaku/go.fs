// Package pkgfs provides a read-only File System based on go import paths.
//
// Essentially it provides a way to select files from a given path, defined in
// terms of a go import path. In your development environment, this will work
// with the real file system and provide access to the files in your GOPATH.
// For production deployment, the included pkgfszip tool will augment a binary
// with a zip file of all the resources it uses, along with its dependencies.
//
// One intended goal with this approach is that the package is still "go get"
// compatible. The binary just needs to be augmented before it can be deployed
// into production.
//
// Another goal is to make the transition from the os package to use this as
// seamless as possible, so where possible the APIs are designed to mimic
// their os counterparts.
package pkgfs

import (
	"go/build"
	"log"
	"os"
	"os/exec"

	"github.com/daaku/go.fs"
	"github.com/daaku/go.fs/emptyfs"
	"github.com/daaku/go.fs/limitfs"
	"github.com/daaku/go.fs/realfs"
	"github.com/daaku/go.fs/zipfs"
)

// Defines a Config that selects files to be made available via a File System.
type Config struct {
	ImportPath string // the import path to use as the root of the File System
	Recursive  bool   // default is not recursive
	Glob       string // optionally limit by a glob pattern
}

// Provides scoped access to a package as a File System. If the currently
// running binary has a zip attached, it will be used, otherwise the GOROOT
// will be used to find the actual files.
func New(c Config) fs.System {
	lc := limitfs.Config{
		Recursive: c.Recursive,
		Glob:      c.Glob,
	}
	var s fs.System
	if exeZipFS == nil {
		s = realfs.New()
		pkg, err := build.Import(c.ImportPath, "", build.FindOnly)
		if err != nil {
			s = emptyfs.New()
		} else {
			lc.Root = pkg.Dir
		}
	} else {
		s = exeZipFS
		lc.Root = c.ImportPath
	}
	return limitfs.New(lc, s)
}

// A singleton zip is expected containing contents from all packages. We open
// this when the process is started and never explicitly close it.
var exeZipFS fs.System

func init() {
	var err error
	exeZipFS, err = openRunningExeAsZip()
	if err != nil {
		log.Println(err)
	}
}

func openRunningExeAsZip() (fs.System, error) {
	name, err := exec.LookPath(os.Args[0])
	if err != nil {
		return nil, err
	}
	s, err := zipfs.Open(name)
	if err != nil {
		return nil, err
	}
	return s, nil
}
