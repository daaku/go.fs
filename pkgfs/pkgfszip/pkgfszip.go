// Command pkgfszip is a build tool that looks for resources used by a package
// (and dependencies) and helps build a zip file with its contents. This allows
// for bundling the used resources into the compiled binary providing for ease
// of distribution.
package main

import (
	"debug/elf"
	"fmt"
	"github.com/daaku/go.atomicfile"
	"github.com/daaku/go.fs/pkgfs/build"
	"github.com/voxelbrain/goptions"
	gobuild "go/build"
	"io"
	"os"
	"path/filepath"
)

// Find the binary path for a command specified as it's import path.
func BinaryPathFromImportPath(importPath, srcDir string) (string, error) {
	pkg, err := gobuild.Import(importPath, srcDir, gobuild.AllowBinary)
	if err != nil {
		return "", err
	}
	return filepath.Join(pkg.BinDir, filepath.Base(pkg.ImportPath)), nil
}

// Find the length of the binary content in an executable. This is useful to
// copy only the binary part and excluding existing zip content appended at the
// end of the file.
func BinaryLength(rda io.ReaderAt) (int64, error) {
	file, err := elf.NewFile(rda)
	if err != nil {
		return 0, err
	}

	var max int64
	for _, sect := range file.Sections {
		if sect.Type == elf.SHT_NOBITS {
			continue
		}

		end := int64(sect.Offset + sect.Size)
		if end > max {
			max = end
		}
	}

	return max, nil
}

// Open an existing file and copy the binary part. The existing file will be
// atomically replaced when this file is closed.
func OpenFile(path string) (io.WriteCloser, error) {
	original, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	end, err := BinaryLength(original)
	if err != nil {
		return nil, err
	}
	if _, err := original.Seek(0, 0); err != nil {
		return nil, err
	}
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	out, err := atomicfile.New(path, fi.Mode())
	if err != nil {
		return nil, err
	}
	if _, err := io.CopyN(out, original, end); err != nil {
		return nil, err
	}
	return out, nil
}

func Main() (err error) {
	options := struct {
		ImportPath string        `goptions:"-p, --package, obligatory, description='package to build'"`
		SrcDir     string        `goptions:"-s, --src-dir, description='src dir for imports'"`
		Verbose    bool          `goptions:"-v, --verbose, description='be verbose'"`
		OutPath    string        `goptions:"-o, --output, description='output file'"`
		Help       goptions.Help `goptions:"-h, --help, description='show this help'"`
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

	build := &build.Build{
		ImportPath: options.ImportPath,
		SrcDir:     options.SrcDir,
		Verbose:    options.Verbose,
		Writer:     out,
	}
	return build.Go()
}

func main() {
	if err := Main(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
