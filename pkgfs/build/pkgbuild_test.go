package build_test

import (
	"archive/zip"
	"bytes"
	"github.com/daaku/go.fs/pkgfs/build"
	"testing"
)

type MemZip struct {
	Buffer *bytes.Buffer
}

func (m *MemZip) Reader() (*zip.Reader, error) {
	return zip.NewReader(bytes.NewReader(m.Buffer.Bytes()), int64(m.Buffer.Len()))
}

func (m *MemZip) Writer() *zip.Writer {
	m.Buffer.Reset()
	return zip.NewWriter(m.Buffer)
}

func NewMemZip() *MemZip {
	m := &MemZip{
		Buffer: new(bytes.Buffer),
	}
	return m
}

func TestSimpleBuild(t *testing.T) {
	t.Parallel()
	memzip := NewMemZip()
	build := &build.Build{
		ImportPath: "github.com/daaku/go.fs/pkgfs/build/test/pkgbuild_test_1",
		Writer:     memzip.Buffer,
	}
	if err := build.Go(); err != nil {
		t.Fatal(err)
	}
	reader, err := memzip.Reader()
	if err != nil {
		t.Fatal(err)
	}
	if l := len(reader.File); l != 1 {
		t.Fatalf("expecting 1 entry in zip got %d", l)
	}
	if reader.File[0].Name != "github.com/daaku/go.fs/pkgfs/build/test/pkgbuild_test_1/main.go" {
		t.Fatalf("did not find expected file, found %s", reader.File[0].Name)
	}
}
