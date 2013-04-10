package pkgbuild

import (
	"archive/zip"
	"fmt"
	"github.com/daaku/go.deepimports"
	"github.com/daaku/go.literalfinder"
	"github.com/daaku/go.pkgfs"
	"io"
	"os"
	"path/filepath"
	"spew"
)

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
	const ref = "github.com/daaku/go.pkgfs.Config"
	b.processed = make(map[string]bool)
	b.zipWriter = zip.NewWriter(b.Writer)
	pkgs, err := deepimports.Find([]string{b.ImportPath}, b.SrcDir)
	if err != nil {
		return err
	}
	for _, pkg := range pkgs {
		if len(pkg.CgoFiles) > 0 {
			if b.Verbose {
				fmt.Printf("skipping %s with cgo files\n", pkg.ImportPath)
			}
			continue
		}
		finder := literalfinder.NewFinder(ref)
		for _, file := range pkg.GoFiles {
			abs := filepath.Join(pkg.SrcRoot, pkg.ImportPath, file)
			if err := b.addSource(finder, abs); err != nil {
				return err
			}
		}

		var configs []*pkgfs.Config
		if err := finder.Find(&configs); err != nil {
			return err
		}
		if len(configs) > 0 {
			spew.Dump(configs)
		}
		for _, config := range configs {
			b.addResource(config)
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

func (b *Build) addSource(finder *literalfinder.Finder, filename string) error {
	if b.Verbose {
		fmt.Printf("Source: %s\n", filename)
	}
	return finder.Add(filename, nil)
}

func (b *Build) addResource(ru *pkgfs.Config) error {
	if b.processed[ru.ImportPath] {
		return nil
	}
	b.processed[ru.ImportPath] = true
	fs := pkgfs.New(*ru)
	root, err := fs.Open("/")
	if err != nil {
		return err
	}
	files, err := root.Readdir(0)
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
		r, err := fs.Open(info.Name())
		if err != nil {
			return err
		}
		if _, err = io.Copy(f, r); err != nil {
			return err
		}
	}
	return nil
}
