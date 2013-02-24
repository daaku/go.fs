package pkgrsrc

import (
	"fmt"
	"github.com/cookieo9/resources-go/v2/resources"
	"go/build"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// A singleton ZipBundle is expected containing contents from all packages.
var zipBundle *resources.ZipBundle

// Provides resources backed by a Package. Note, packages _do not_ recurse and
// only provide the top level files. The APIs behave much like their
// counterparts in the os package.
type PackageResources interface {
	// Open a named file for reading.
	Open(path string) (io.ReadCloser, error)

	// Stat returns a FileInfo describing the named file.
	Stat(name string) (fi os.FileInfo, err error)

	// Get the list of sorted FileInfos for all the files in the package.
	ReadDir() ([]os.FileInfo, error)

	// IsNotExist returns whether the error is known to report that a file does not
	// exist.
	IsNotExist(err error) bool
}

// Provides scoped access to a package as raw resources. If the currently
// running binary has a zip attached, it will be used, otherwise the GOROOT
// will be used to find the actual files.
func New(path string) PackageResources {
	if zipBundle == nil {
		return &dir{path, ""}
	}
	return &zip{path}
}

type dir struct {
	path     string
	realPath string
}

func (d *dir) Open(path string) (io.ReadCloser, error) {
	root, err := d.root()
	if err != nil {
		return nil, err
	}
	return os.Open(filepath.Join(root, path))
}

func (d *dir) Stat(path string) (os.FileInfo, error) {
	root, err := d.root()
	if err != nil {
		return nil, err
	}
	return os.Stat(filepath.Join(root, path))
}

func (d *dir) ReadDir() ([]os.FileInfo, error) {
	root, err := d.root()
	if err != nil {
		return nil, err
	}
	fis, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var finew []os.FileInfo
	for _, fi := range fis {
		if fi.Mode()&os.ModeType != 0 {
			continue
		}
		finew = append(finew, fi)
	}
	return finew, nil
}

func (d *dir) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func (d *dir) root() (string, error) {
	if d.realPath == "" {
		pkg, err := build.Import(d.path, "", build.FindOnly)
		if err != nil {
			return "", fmt.Errorf("package provider import path %s not found", d.path)
		}
		d.realPath = pkg.Dir
	}
	return d.realPath, nil
}

type zip struct {
	root string
}

func (z *zip) Open(path string) (io.ReadCloser, error) {
	return zipBundle.Open(filepath.Join(z.root, path))
}

func (z *zip) Stat(path string) (os.FileInfo, error) {
	r, err := zipBundle.Find(filepath.Join(z.root, path))
	if err != nil {
		return nil, err
	}
	return r.Stat()
}

func (z *zip) ReadDir() ([]os.FileInfo, error) {
	rs, err := zipBundle.Glob(filepath.Join(z.root, "*"))
	if err != nil {
		return nil, err
	}
	var fis []os.FileInfo
	for _, r := range rs {
		fi, err := r.Stat()
		if err != nil {
			return nil, err
		}
		fis = append(fis, fi)
	}
	return fis, err
}

func (z *zip) IsNotExist(err error) bool {
	return err == resources.ErrNotFound
}

func init() {
	if p, _ := exec.LookPath(os.Args[0]); p != "" {
		zipBundle, _ = resources.OpenZip(p)
	}
}
