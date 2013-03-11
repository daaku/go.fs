// Package pkgfs provides a read-only FileSystem based around go import
// paths.
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
	"archive/zip"
	"errors"
	"fmt"
	"go/build"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

var errInvalidCharacterInPath = errors.New("invalid character in file path")

// A singleton zip is expected containing contents from all packages.
var bundle *zip.Reader

// A read-only File as returned by a FileSystem's Open method.
type File interface {
	io.ReadCloser
	Stat() (os.FileInfo, error)
	Readdir(count int) ([]os.FileInfo, error)
}

// A FileSystem implements access to a collection of named files. The elements
// in a file path are separated by slash ('/', U+002F) characters, regardless
// of host operating system convention.
type FileSystem interface {
	// Open a named file for reading.
	Open(name string) (File, error)

	// IsNotExist returns whether the error is known to report that a file does
	// not exist.
	IsNotExist(err error) bool
}

// Defines a Config that selects files to be made available via a FileSystem.
type Config struct {
	ImportPath string // the import path to use as the root of the FileSystem
	Recursive  bool   // default is not recursive
	Glob       string // optionally limit by a glob pattern
}

// Provides scoped access to a package as a FileSystem. If the currently
// running binary has a zip attached, it will be used, otherwise the GOROOT
// will be used to find the actual files.
func New(c Config) FileSystem {
	if bundle == nil {
		return &dirFS{
			path:      c.ImportPath,
			recursive: c.Recursive,
			glob:      c.Glob,
		}
	}
	return &zipFS{
		root:      c.ImportPath,
		bundle:    bundle,
		recursive: c.Recursive,
		glob:      c.Glob,
	}
}

type errNotFound string

func (e errNotFound) Error() string {
	return fmt.Sprintf("file not found: %s", string(e))
}

type errNotIncluded string

func (e errNotIncluded) Error() string {
	return fmt.Sprintf("file not included: %s", string(e))
}

func isNotExist(err error) bool {
	_, ok := err.(errNotFound)
	if ok {
		return true
	}
	_, ok = err.(errNotIncluded)
	return ok
}

func cleanName(name string, recursive bool, glob string) (string, error) {
	if !recursive && len(filepath.SplitList(name)) > 1 {
		return "", errNotIncluded(name)
	}
	if filepath.Separator != '/' &&
		strings.IndexRune(name, filepath.Separator) >= 0 ||
		strings.Contains(name, "\x00") {
		return "", errInvalidCharacterInPath
	}
	clean := filepath.FromSlash(path.Clean("/" + name))
	if glob != "" {
		match, err := path.Match(glob, name)
		if err != nil {
			return "", err
		}
		if !match {
			return "", errNotIncluded(name)
		}
	}
	return clean, nil
}

type dirFS struct {
	path      string
	realPath  string
	recursive bool
	glob      string
}

func (d *dirFS) Open(name string) (File, error) {
	root, err := d.root()
	if err != nil {
		return nil, err
	}
	clean, err := cleanName(name, d.recursive, d.glob)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(filepath.Join(root, clean))
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (d *dirFS) IsNotExist(err error) bool {
	return os.IsNotExist(err) || isNotExist(err)
}

func (d *dirFS) root() (string, error) {
	if d.realPath == "" {
		pkg, err := build.Import(d.path, "", build.FindOnly)
		if err != nil {
			return "", fmt.Errorf("filesystem at import path %s not found", d.path)
		}
		d.realPath = pkg.Dir
	}
	return d.realPath, nil
}

type zipFile struct {
	io.ReadCloser
	*zip.File
}

func (z *zipFile) Readdir(count int) ([]os.FileInfo, error) {
	//TODO
	return nil, nil
}

func (z *zipFile) Stat() (os.FileInfo, error) {
	return z.FileInfo(), nil
}

type zipFS struct {
	root      string
	bundle    *zip.Reader
	recursive bool
	glob      string
}

func (z *zipFS) Open(name string) (File, error) {
	clean, err := cleanName(name, z.recursive, z.glob)
	if err != nil {
		return nil, err
	}
	for _, f := range z.bundle.File {
		if f.Name == clean {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			return &zipFile{
				ReadCloser: rc,
				File:       f,
			}, nil
		}
	}
	return nil, errNotFound(name)
}

func (z *zipFS) IsNotExist(err error) bool {
	return isNotExist(err)
}

func init() {
	if p, _ := exec.LookPath(os.Args[0]); p != "" {
		rc, err := zip.OpenReader(p)
		if err != nil {
			fmt.Println(err)
			return
		}
		bundle = &rc.Reader
	}
}
