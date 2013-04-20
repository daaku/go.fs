package memfs_test

import (
	"os"
	"testing"
	"time"

	"github.com/daaku/go.fs/memfs"
)

func TestFileInfoFromExisting(t *testing.T) {
	t.Parallel()
	const name = "foo"
	fi := memfs.NewFileInfoFromExisting(
		memfs.NewFileInfo(memfs.FileInfo{
			Name: name,
		}),
	)
	if fi.Name() != name {
		t.Fatal("did not find expected name")
	}
}

func TestFileInfoSetSys(t *testing.T) {
	t.Parallel()
	fi := memfs.NewFileInfo(memfs.FileInfo{})
	var s interface{} = 42
	fi.SetSys(s)
	if fi.Sys() != s {
		t.Fatal("did not find expected sys")
	}
}

func TestFileInfoIsDir(t *testing.T) {
	t.Parallel()
	fi := memfs.NewFileInfo(memfs.FileInfo{
		Mode: os.ModeDir,
	})
	if !fi.IsDir() {
		t.Fatal("was expecting dir")
	}
}

func TestFileInfoSetModTime(t *testing.T) {
	t.Parallel()
	now := time.Now()
	fi := memfs.NewFileInfo(memfs.FileInfo{
		ModTime: now,
	})
	if !fi.ModTime().Equal(now) {
		t.Fatalf("was expecting time %s got %s", now, fi.ModTime())
	}
	other := now.AddDate(10, 0, 0)
	fi.SetModTime(other)
	if !fi.ModTime().Equal(other) {
		t.Fatalf("was expecting time %s got %s", other, fi.ModTime())
	}
}
