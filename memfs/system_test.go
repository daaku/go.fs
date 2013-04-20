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
