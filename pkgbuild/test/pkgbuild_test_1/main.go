package main

import (
	"github.com/daaku/go.pkgrsrc/pkgrsrc"
)

var FS = pkgrsrc.New(pkgrsrc.Config{
	ImportPath: "github.com/daaku/go.pkgrsrc/pkgbuild/test/pkgbuild_test_1",
	Recursive:  true,
})

func main() {
}
