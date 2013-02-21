// Command pkgrsrc is a build tool that looks for resources used by a package
// (and dependencies) and helps build a zip file with its contents. This allows
// for bundling the used resources into the compiled binary providing for ease
// of distribution.
package main

import (
	"fmt"
	"github.com/daaku/go.pkgrsrc/pkgrsrc"
	"github.com/voxelbrain/goptions"
	"go/build"
	"io"
	"os"
	"path/filepath"
)

func BinaryPathFromImportPath(importPath, srcDir string) (string, error) {
	pkg, err := build.Import(importPath, srcDir, build.AllowBinary)
	if err != nil {
		return "", err
	}
	return filepath.Join(pkg.BinDir, filepath.Base(pkg.ImportPath)), nil
}

func OpenFile(path string) (io.WriteCloser, error) {
	// TODO atomic rename
	// TODO copy existing binary and append
	// TODO handle binary with existing zip content
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, fi.Mode())
	if err != nil {
		return nil, err
	}
	return file, nil
}

func Main() (err error) {
	options := struct {
		ImportPath string `goptions:"-p, --package, obligatory, description='package to build'"`
		SrcDir     string `goptions:"-s, --src-dir, description='src dir for imports'"`
		Verbose    bool   `goptions:"-v, --verbose, description='be verbose'"`
		OutPath    string `goptions:"-o, --output, description='output file'"`
	}{}
	goptions.ParseAndFail(&options)

	if options.OutPath == "" {
		options.OutPath, err = BinaryPathFromImportPath(
			options.ImportPath, options.SrcDir)
		if err != nil {
			return err
		}
	}

	out, err := OpenFile(options.OutPath)
	if err != nil {
		return err
	}
	defer out.Close()

	build := &pkgrsrc.Build{
		ImportPath: options.ImportPath,
		SrcDir:     options.SrcDir,
		Verbose:    options.Verbose,
		Writer:     out,
	}
	return build.Go()
}

func main() {
	if err := Main(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}
}
