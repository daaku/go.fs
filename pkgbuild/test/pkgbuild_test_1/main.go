package main

import (
	"github.com/daaku/go.commonjs"
	"github.com/daaku/go.pkgrsrc/pkgrsrc"
)

func Foo() commonjs.Provider {
	return commonjs.NewFileSystemProvider(
		pkgrsrc.New(pkgrsrc.Config{
			ImportPath: "github.com/daaku/go.pkgrsrc/pkgbuild/test/pkgbuild_test_1",
			Recursive:  true,
		}))
}

func main() {
}
