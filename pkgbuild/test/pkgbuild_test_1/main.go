package main

import (
	"github.com/daaku/go.commonjs"
	"github.com/daaku/go.pkgrsrc/pkgrsrc"
)

func Foo() commonjs.Provider {
	return commonjs.NewFileSystemProvider(
		pkgrsrc.New("github.com/daaku/go.pkgrsrc/pkgbuild/test/pkgbuild_test_1"))
}

func main() {
}
