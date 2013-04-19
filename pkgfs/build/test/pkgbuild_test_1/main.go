package main

import (
	"github.com/daaku/go.fs/pkgfs"
)

var FS = pkgfs.New(pkgfs.Config{
	ImportPath: "github.com/daaku/go.pkgfs/pkgbuild/test/pkgbuild_test_1",
	Recursive:  true,
})

func main() {
}
