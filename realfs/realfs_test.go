package realfs_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/daaku/go.fs/realfs"
)

func TestOpenExisting(t *testing.T) {
	s := realfs.New()
	_, err := s.Open("/etc/passwd")
	if err != nil {
		t.Fatal(err)
	}
}

func TestOpenNotExist(t *testing.T) {
	s := realfs.New()
	_, err := s.Open("/foo/bar/baz/boom")
	if err == nil {
		t.Fatal("was expecting error")
	}
	if !s.IsNotExist(err) {
		t.Fatal("was expecting is not exist error")
	}
}

func TestGetUGID(t *testing.T) {
	tf, err := ioutil.TempFile("", "realfs_test")
	if err != nil {
		t.Fatal(err)
	}
	name := tf.Name()
	defer os.Remove(name)
	s := realfs.New()
	f, err := s.Open(name)
	if err != nil {
		t.Fatal(err)
	}
	expectedUID, expectedGID := os.Getuid(), os.Getgid()
	actualGID, err := f.OwnerGID()
	if err != nil {
		t.Fatal(err)
	}
	if actualGID != expectedGID {
		t.Fatal("did not find expected gid")
	}
	actualUID, err := f.OwnerUID()
	if err != nil {
		t.Fatal(err)
	}
	if actualUID != expectedUID {
		t.Fatal("did not find expected uid")
	}
}

func TestGetUGIDError(t *testing.T) {
	tf, err := ioutil.TempFile("", "realfs_test")
	if err != nil {
		t.Fatal(err)
	}
	name := tf.Name()
	s := realfs.New()
	f, err := s.Open(name)
	if err != nil {
		t.Fatal(err)
	}
	tf.Close()
	f.Close()
	os.Remove(name)
	if _, err := f.OwnerGID(); err == nil {
		t.Fatal("was expecting error")
	}
	if _, err := f.OwnerUID(); err == nil {
		t.Fatal("was expecting error")
	}
}
