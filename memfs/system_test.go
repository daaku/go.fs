package memfs_test

import (
	"os"
	"testing"
	"time"

	"github.com/daaku/go.fs"
	"github.com/daaku/go.fs/memfs"
)

func TestSimpleSystem(t *testing.T) {
	t.Parallel()
	const name = "foo"
	const data = "bar"
	createdFile := memfs.NewFile(name, os.FileMode(666), time.Now(), []byte(data))
	s := memfs.NewSystem(map[string]fs.File{name: createdFile})
	openedFile, err := s.Open(name)
	if err != nil {
		t.Fatal(err)
	}
	if openedFile != createdFile {
		t.Fatal("did not find expected file")
	}
}

func TestSystemNotFound(t *testing.T) {
	t.Parallel()
	s := memfs.NewSystem(nil)
	_, err := s.Open("foo")
	if err == nil {
		t.Fatal("was expecting error")
	}
	if !s.IsNotExist(err) {
		t.Fatal("was expecting is not exist error")
	}
}

func TestSystemWithFiles(t *testing.T) {
	t.Parallel()
	f1 := memfs.NewFile("foo", os.FileMode(666), time.Now(), nil)
	f2 := memfs.NewFile("bar", os.FileMode(666), time.Now(), nil)
	s := memfs.NewWithFiles(map[string]fs.File{
		"d/foo": f1,
		"d/bar": f2,
	})
	a1, err := s.Open("d/foo")
	if err != nil {
		t.Fatal(err)
	}
	if a1 != f1 {
		t.Fatal("did not find expected file")
	}
	d, err := s.Open("d")
	if err != nil {
		t.Fatal(err)
	}
	some, err := d.Readdirnames(0)
	if len(some) != 2 {
		t.Fatal("was expecting 2 names")
	}
	if actual := some[1]; actual != "bar" {
		t.Fatal("was expecting bar")
	}
}

func TestSystemWithFilesClosedFile(t *testing.T) {
	t.Parallel()
	f1 := memfs.NewFile("foo", os.FileMode(666), time.Now(), nil)
	f1.Close()
	s := memfs.NewWithFiles(map[string]fs.File{
		"d/foo": f1,
	})
	_, err := s.Open("d/foo")
	if err == nil {
		t.Fatal("was expecting error")
	}
}

func TestSystemWithFilesIncorrectStructure(t *testing.T) {
	t.Parallel()
	f1 := memfs.NewFile("foo", os.FileMode(666), time.Now(), nil)
	f2 := memfs.NewFile("bar", os.FileMode(666), time.Now(), nil)
	s := memfs.NewWithFiles(map[string]fs.File{
		"d":     f2,
		"d/foo": f1,
	})
	_, err := s.Open("d/foo")
	if err == nil {
		t.Fatal("was expecting error")
	}
}

func TestSystemWithFilesTwoLevels(t *testing.T) {
	t.Parallel()
	f1 := memfs.NewFile("foo", os.FileMode(666), time.Now(), nil)
	f2 := memfs.NewFile("bar", os.FileMode(666), time.Now(), nil)
	s := memfs.NewWithFiles(map[string]fs.File{
		"d/e/foo": f1,
		"d/e/bar": f2,
	})
	a1, err := s.Open("d/e/foo")
	if err != nil {
		t.Fatal(err)
	}
	if a1 != f1 {
		t.Fatal("did not find expected file")
	}
	d, err := s.Open("d")
	if err != nil {
		t.Fatal(err)
	}
	some, err := d.Readdirnames(0)
	if len(some) != 1 {
		t.Fatal("was expecting 1 name")
	}
	if actual := some[0]; actual != "e" {
		t.Fatal("was expecting bar")
	}
}
