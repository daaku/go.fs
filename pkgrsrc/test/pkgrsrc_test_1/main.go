package main

import (
	"github.com/daaku/go.commonjs"
)

func Foo() commonjs.Provider {
	return commonjs.NewPackageProvider("github.com/daaku/go.pkgrsrc/pkgrsrc/test/pkgrsrc_test_1")
}