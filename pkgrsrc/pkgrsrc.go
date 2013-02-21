package pkgrsrc

import (
	"archive/zip"
	"fmt"
	"github.com/daaku/go.deepimports"
	"go/build"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

var reFunCall = regexp.MustCompile(`NewPackageProvider\(['"](.+?)['"]\)`)

type ResourceUsage struct {
	ImportPath string
}

// Parses source for resource usage and returns a list of import paths that are
// referenced.
func ParseResourceUsage(content []byte) ([]*ResourceUsage, error) {
	calls := reFunCall.FindAllSubmatch(content, -1)
	l := make([]*ResourceUsage, len(calls))
	for ix, dep := range calls {
		l[ix] = &ResourceUsage{ImportPath: string(dep[1])}
	}
	return l, nil
}

// Builds a zip file for a specific package by import path.
type Build struct {
	ImportPath string // the package to build for
	SrcDir     string // src dir for finding packages
	Verbose    bool
	Writer     io.Writer
	zipWriter  *zip.Writer
	processed  map[string]bool
}

// Build and write the zip file.
func (b *Build) Go() error {
	b.processed = make(map[string]bool)
	b.zipWriter = zip.NewWriter(b.Writer)
	pkgs, err := deepimports.Find([]string{b.ImportPath}, b.SrcDir)
	if err != nil {
		return err
	}
	for _, pkg := range pkgs {
		for _, file := range pkg.GoFiles {
			abs := filepath.Join(pkg.SrcRoot, pkg.ImportPath, file)
			if err := b.parseAndAddResources(abs); err != nil {
				return err
			}
		}
		for _, file := range pkg.CgoFiles {
			abs := filepath.Join(pkg.SrcRoot, pkg.ImportPath, file)
			if err := b.parseAndAddResources(abs); err != nil {
				return err
			}
		}
	}

	if b.Verbose {
		fmt.Println("closing zip file")
	}
	if err := b.zipWriter.Close(); err != nil {
		return err
	}
	return nil
}

func (b *Build) parseAndAddResources(path string) error {
	if b.Verbose {
		fmt.Printf("Source: %s\n", path)
	}
	rus, err := b.parseResourceUsage(path)
	if err != nil {
		return err
	}
	for _, ru := range rus {
		if err := b.addResource(ru); err != nil {
			return err
		}
	}
	return nil
}

func (b *Build) addResource(ru *ResourceUsage) error {
	if b.processed[ru.ImportPath] {
		return nil
	}
	b.processed[ru.ImportPath] = true
	pkg, err := build.Import(ru.ImportPath, b.SrcDir, build.AllowBinary)
	if err != nil {
		return err
	}
	rootAbs := filepath.Join(pkg.SrcRoot, pkg.ImportPath)
	if b.Verbose {
		fmt.Printf("Package: [%s]/%s\n", pkg.SrcRoot, pkg.ImportPath)
	}
	rootDir, err := os.Open(rootAbs)
	if err != nil {
		return err
	}
	files, err := rootDir.Readdir(0)
	if err != nil {
		return err
	}
	for _, info := range files {
		zabs := filepath.Join(ru.ImportPath, info.Name())
		if info.Mode()&os.ModeType != 0 {
			if b.Verbose {
				fmt.Printf("Skipped Resource: %s\n", zabs)
			}
			continue
		}
		if b.Verbose {
			fmt.Printf("Resource: %s\n", zabs)
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = zabs
		f, err := b.zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		rabs := filepath.Join(rootAbs, info.Name())
		r, err := os.Open(rabs)
		if err != nil {
			return err
		}
		if _, err = io.Copy(f, r); err != nil {
			return err
		}
	}
	return nil
}

func (b *Build) parseResourceUsage(path string) ([]*ResourceUsage, error) {
	r, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseResourceUsage(r)
}
