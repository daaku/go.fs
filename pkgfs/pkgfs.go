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
	"debug/elf"
	"errors"
	"fmt"
	"go/build"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/daaku/go.zipexe"
)

var errInvalidCharacterInPath = errors.New("invalid character in file path")

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
	base := baseFS{
		recursive: c.Recursive,
		glob:      c.Glob,
	}
	if bundle == nil {
		return &dirFS{
			baseFS: base,
			path:   c.ImportPath,
		}
	}
	return &zipFS{
		baseFS: base,
		root:   c.ImportPath,
		bundle: bundle,
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

type limitedDir struct {
	File
	FileSystem FileSystem
	glob       string
	recusrive  bool
}

func (l limitedDir) Readdir(count int) (fis []os.FileInfo, err error) {
	if count <= 0 {
		raw, err := l.File.Readdir(count)
		fis, errF := l.filter(raw)
		if err == nil && errF != nil {
			err = errF
		}
		return fis, err
	}

	pending := count
	for pending > 0 {
		raw, err := l.File.Readdir(pending)
		if len(raw) == 0 {
			return fis, err
		}
		filtered, errF := l.filter(raw)
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

func (l limitedDir) filter(given []os.FileInfo) (final []os.FileInfo, err error) {
	// FIXME
	for _, fi := range given {
		if fi.IsDir() {
			if !l.recusrive {
				continue
			}
			//cfs, err := 1, nil
		}
		final = append(final, fi)
	}
	return final, nil
}

type baseFS struct {
	recursive bool
	glob      string
}

func (b baseFS) cleanName(name string) (string, error) {
	if filepath.Separator != '/' &&
		strings.IndexRune(name, filepath.Separator) >= 0 ||
		strings.Contains(name, "\x00") {
		return "", errInvalidCharacterInPath
	}
	clean := filepath.FromSlash(path.Clean("/" + name))
	if !b.recursive && len(filepath.SplitList(clean)) > 1 {
		return "", errNotIncluded(name)
	}
	return clean, nil
}

func (b baseFS) isIncluded(name string) (bool, error) {
	if b.glob == "" {
		return true, nil
	}
	match, err := path.Match(b.glob, name)
	if err != nil {
		return false, err
	}
	return match, nil
}

type dirFS struct {
	baseFS
	path     string
	realPath string
}

func (d *dirFS) Open(name string) (File, error) {
	root, err := d.root()
	if err != nil {
		return nil, err
	}
	clean, err := d.cleanName(name)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(filepath.Join(root, clean))
	if err != nil {
		return nil, err
	}
	included, err := d.isIncluded(clean)
	if err != nil {
		return nil, err
	}
	if !included {
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
	baseFS
	root   string
	bundle *zip.Reader
}

func (z *zipFS) Open(name string) (File, error) {
	clean, err := z.cleanName(name)
	if err != nil {
		return nil, err
	}
	fmt.Println(clean)
	if len(z.bundle.File) > 0 {
		fmt.Println(len(z.bundle.File))
	}
	for _, f := range z.bundle.File {
		fmt.Println(f.Name)
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

// A singleton zip is expected containing contents from all packages. We open
// this when the process is started and never explicitly close it.
var exeZip *zip.Reader

func init() {
	var err error
	exeZip, err = openBundle()
	if err != nil {
		fmt.Println(err)
	}
}

func openBundle() (*zip.Reader, error) {
	name, err := exec.LookPath(os.Args[0])
	if err != nil {
		return nil, err
	}
	zr, err := zipexe.Open(file)
	if err != nil {
		return nil, err
	}
	return zr, nil
}
