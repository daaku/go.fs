package main

import (
	"github.com/daaku/go.commonjs"
)

func Foo() commonjs.Provider {
	return commonjs.NewPackageResourceProvider("github.com/daaku/go.pkgrsrc/pkgrsrc/test/pkgrsrc_test_1")
}

func main() {
}
