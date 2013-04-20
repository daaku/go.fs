package emptyfs_test

import (
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
