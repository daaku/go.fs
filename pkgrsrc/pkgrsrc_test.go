package pkgrsrc_test

import (
	"archive/zip"
	"bytes"
	"github.com/daaku/go.pkgrsrc/pkgrsrc"
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
	memzip := NewMemZip()
	build := &pkgrsrc.Build{
		ImportPath: "github.com/daaku/go.pkgrsrc/pkgrsrc/test/pkgrsrc_test_1",
		Writer:     memzip.Buffer,
	}
	if err := build.Go(); err != nil {
		t.Fatal(err)
	}
	reader, err := memzip.Reader()
	if err != nil {
		t.Fatal(err)
	}
	if len(reader.File) != 1 {
		t.Fatal("expecting 1 entry in zip")
	}
	if reader.File[0].Name != "github.com/daaku/go.pkgrsrc/pkgrsrc/test/pkgrsrc_test_1/main.go" {
		t.Fatalf("did not find expected file, found %s", reader.File[0].Name)
	}
}

func TestParseResourceUsage(t *testing.T) {
	content := []byte(`NewPackageProvider("foo")`)
	rus, err := pkgrsrc.ParseResourceUsage(content)
	if err != nil {
		t.Fatal(err)
	}
	if len(rus) != 1 {
		t.Fatal("was expecting 1 resource usage")
	}
	if rus[0].ImportPath != "foo" {
		t.Fatalf("was expecting foo resource usage but got %s", rus[0].ImportPath)
	}
}
