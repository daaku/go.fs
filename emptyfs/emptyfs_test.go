package emptyfs_test

import (
	"errors"
	"testing"

	"github.com/daaku/go.fs/emptyfs"
)

func TestAlwaysEmpty(t *testing.T) {
	t.Parallel()
	s := emptyfs.New()
	_, err := s.Open("foo")
	if err == nil {
		t.Fatal("expecting error")
	}
	if !s.IsNotExist(err) {
		t.Fatal("expecting is not exist error")
	}
}

func TestFixedError(t *testing.T) {
	t.Parallel()
	e := errors.New("t")
	s := emptyfs.NewWithError(e)
	_, err := s.Open("foo")
	if err == nil {
		t.Fatal("expecting error")
	}
	if e != err {
		t.Fatal("expecting configured error")
	}
	if !s.IsNotExist(err) {
		t.Fatal("expecting is not exist error")
	}
}
