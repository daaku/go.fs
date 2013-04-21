package memfs_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/daaku/go.fs/memfs"
)

var (
	dMode = os.FileMode(666)
	dTime = time.Now()
)

func TestFileChmod(t *testing.T) {
	t.Parallel()
	f := memfs.NewFile("foo", dMode, dTime, nil)
	stat, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if stat.Mode() != dMode {
		t.Fatal("did not find original mode")
	}
	nMode := os.FileMode(777)
	err = f.Chmod(nMode)
	if err != nil {
		t.Fatal(err)
	}
	stat, err = f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if stat.Mode() != nMode {
		t.Fatal("did not find new mode")
	}
}

func TestFileDefaultOwner(t *testing.T) {
	t.Parallel()
	f := memfs.NewFile("foo", dMode, dTime, nil)
	gid, err := f.OwnerGID()
	if err != nil {
		t.Fatal(err)
	}
	if gid != 0 {
		t.Fatal("expected 0 gid")
	}
	uid, err := f.OwnerUID()
	if err != nil {
		t.Fatal(err)
	}
	if uid != 0 {
		t.Fatal("expected 0 uid")
	}
}

func TestFileChown(t *testing.T) {
	t.Parallel()
	f := memfs.NewFile("foo", dMode, dTime, nil)
	expectedUID, expectedGID := 1, 2
	err := f.Chown(expectedUID, expectedGID)
	if err != nil {
		t.Fatal(err)
	}
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

func TestFileNameOnCreate(t *testing.T) {
	t.Parallel()
	f := memfs.NewFile("foo", dMode, dTime, nil)
	stat, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if stat.Name() != "foo" {
		t.Fatal("did not find expected name")
	}
}

func TestFileNameOnSet(t *testing.T) {
	t.Parallel()
	f := memfs.NewFile("foo", dMode, dTime, nil)
	name := "bar/baz"
	f.SetName(name)
	if f.Name() != name {
		t.Fatal("did not find expected name")
	}
	stat, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if stat.Name() != "baz" {
		t.Fatal("did not find expected name")
	}
}

func TestFileReadLess(t *testing.T) {
	t.Parallel()
	f := memfs.NewFile("foo", dMode, dTime, []byte("ab"))
	b := make([]byte, 1)
	n, err := f.Read(b)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatal("did not find expected count")
	}
	if b[0] != 'a' {
		t.Fatal("did not find expected byte")
	}
	n, err = f.Read(b)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatal("did not find expected count")
	}
	if b[0] != 'b' {
		t.Fatal("did not find expected byte")
	}
	n, err = f.Read(nil)
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatal("did not find expected count")
	}
	n, err = f.Read(b)
	if err != io.EOF {
		t.Fatalf("was expecting EOF %s", err)
	}
	if n != 0 {
		t.Fatal("did not find expected count")
	}
	if b[0] != 'b' {
		t.Fatal("did not find expected byte")
	}
}

func TestFileReadReset(t *testing.T) {
	t.Parallel()
	f := memfs.NewFile("foo", dMode, dTime, []byte("ab"))
	b := make([]byte, 1)
	n, err := f.Read(b)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatal("did not find expected count")
	}
	if b[0] != 'a' {
		t.Fatal("did not find expected byte")
	}
	f.Reset()
	n, err = f.Read(b)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatal("did not find expected count")
	}
	if b[0] != 'a' {
		t.Fatal("did not find expected byte")
	}
}

func TestFileReadMore(t *testing.T) {
	t.Parallel()
	f := memfs.NewFile("foo", dMode, dTime, []byte("ab"))
	b := make([]byte, 10)
	n, err := f.Read(b)
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatal("did not find expected count")
	}
	if string(b[:2]) != "ab" {
		t.Fatalf("did not find expected bytes, found: %v", b)
	}
}

func TestFileReadAt(t *testing.T) {
	t.Parallel()
	f := memfs.NewFile("foo", dMode, dTime, []byte("ab"))
	b := make([]byte, 10)
	n, err := f.ReadAt(b, 1)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatal("did not find expected count")
	}
	if b[0] != 'b' {
		t.Fatal("did not find expected byte")
	}
}

func TestFileSeek(t *testing.T) {
	t.Parallel()
	f := memfs.NewFile("foo", dMode, dTime, []byte("123456789"))
	n, err := f.Seek(2, os.SEEK_CUR)
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("was expecting 2 got %d", n)
	}
	n, err = f.Seek(2, os.SEEK_SET)
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("was expecting 2 got %d", n)
	}
	n, err = f.Seek(7, os.SEEK_END)
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("was expecting 2 got %d", n)
	}
	n, err = f.Seek(7, os.SEEK_CUR)
	if err != nil {
		t.Fatal(err)
	}
	if n != 9 {
		t.Fatalf("was expecting 9 got %d", n)
	}
	_, err = f.Seek(1, os.SEEK_CUR)
	if err == nil || !strings.Contains(err.Error(), "out of range") {
		t.Fatalf("was expecting out of range error: %s", err)
	}
	_, err = f.Seek(-1, os.SEEK_SET)
	if err == nil || !strings.Contains(err.Error(), "out of range") {
		t.Fatalf("was expecting out of range error: %s", err)
	}
	_, err = f.Seek(-1, 99)
	if err == nil || !strings.Contains(err.Error(), "whence invalid") {
		t.Fatalf("was expecting whence error: %s", err)
	}
}

func TestFileTruncate(t *testing.T) {
	t.Parallel()
	in := []byte("123456789")
	f := memfs.NewFile("foo", dMode, dTime, in)
	stat, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if stat.Size() != 9 {
		t.Fatal("did not find expected size in FileInfo")
	}
	err = f.Truncate(4)
	if err != nil {
		t.Fatal(err)
	}
	stat, err = f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if stat.Size() != 4 {
		t.Fatal("did not find expected size in FileInfo")
	}
	out, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(in[:4], out) {
		t.Fatalf("did not get the same bytes: in: %v out: %v", in[:4], out)
	}
	err = f.Truncate(0)
	if err != nil {
		t.Fatal(err)
	}
	stat, err = f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if stat.Size() != 0 {
		t.Fatal("did not find expected size in FileInfo")
	}
	if !strings.Contains(f.Truncate(1).Error(), "out of range") {
		t.Fatal("was expecting out of range")
	}
	if !strings.Contains(f.Truncate(-42).Error(), "out of range") {
		t.Fatal("was expecting out of range")
	}
}

func TestFileTruncateGrowToEnsureReuse(t *testing.T) {
	t.Parallel()
	in := []byte("123456789")
	f := memfs.NewFile("foo", dMode, dTime, in)
	err := f.Truncate(4)
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.Seek(0, os.SEEK_END)
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString("abc")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.Seek(0, os.SEEK_SET)
	if err != nil {
		t.Fatal(err)
	}
	out, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal([]byte("1234abc"), out) {
		t.Fatalf("did not get the same bytes: %s", out)
	}
}

func TestFileWriteAt(t *testing.T) {
	t.Parallel()
	f := memfs.NewFile("foo", dMode, dTime, nil)
	in := []byte("foo bar")
	n, err := f.WriteAt(in, 0)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(in) {
		t.Fatal("did not find expected count")
	}
	out, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(in, out) {
		t.Fatalf("did not get the same bytes: %s", out)
	}
}

func TestFileWriteAtOutOfRange(t *testing.T) {
	t.Parallel()
	f := memfs.NewFile("foo", dMode, dTime, nil)
	_, err := f.WriteAt([]byte("a"), 1)
	if err == nil || !strings.Contains(err.Error(), "out of range") {
		t.Fatalf("was expecting already closed, got %s", err)
	}
}

func TestFileWriteToNil(t *testing.T) {
	t.Parallel()
	f := memfs.NewFile("foo", dMode, dTime, nil)
	in := []byte("foo bar")
	n, err := f.Write(in)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(in) {
		t.Fatal("did not find expected count")
	}
	stat, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if stat.Size() != int64(n) {
		t.Fatal("did not find expected size in FileInfo")
	}
}

func TestFileClosed(t *testing.T) {
	t.Parallel()
	f := memfs.NewFile("foo", dMode, dTime, nil)
	assertClosed := func(err error) {
		if !strings.Contains(err.Error(), "already closed") {
			t.Fatalf("was expecting already closed, got %s", err)
		}
	}
	err := f.Close()
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.Read(nil)
	assertClosed(err)
	_, err = f.ReadAt(nil, 2)
	assertClosed(err)
	_, err = f.Seek(2, os.SEEK_SET)
	assertClosed(err)
	_, err = f.Stat()
	assertClosed(err)
	err = f.Sync()
	assertClosed(err)
	err = f.Truncate(5)
	assertClosed(err)
	_, err = f.Write(nil)
	assertClosed(err)
	_, err = f.WriteAt(nil, 5)
	assertClosed(err)
	_, err = f.WriteString("")
	assertClosed(err)
}

func TestFileSync(t *testing.T) {
	t.Parallel()
	if memfs.NewFile("foo", dMode, dTime, nil).Sync() != nil {
		t.Fatal("failed sync")
	}
}

func TestFileWithDirOperations(t *testing.T) {
	t.Parallel()
	f := memfs.NewFile("foo", dMode, dTime, nil)
	assertNotDir := func(err error) {
		if !strings.Contains(err.Error(), "is not a directory") {
			t.Fatalf("was expecting is not a directory, got %s", err)
		}
	}
	_, err := f.Readdir(0)
	assertNotDir(err)
	_, err = f.Readdirnames(0)
	assertNotDir(err)
	err = f.SetDirInfos(nil)
	assertNotDir(err)
	err = f.AddDirInfo(nil)
	assertNotDir(err)
}

func TestDirWithFileOperations(t *testing.T) {
	t.Parallel()
	d := memfs.NewDir("foo", dMode, dTime, nil)
	assertIsDir := func(err error) {
		if !strings.Contains(err.Error(), "is a directory") {
			t.Fatalf("was expecting is a directory, got %s", err)
		}
	}
	_, err := d.Read(nil)
	assertIsDir(err)
	_, err = d.ReadAt(nil, 2)
	assertIsDir(err)
	_, err = d.Seek(2, os.SEEK_SET)
	assertIsDir(err)
	err = d.Truncate(5)
	assertIsDir(err)
	_, err = d.Write(nil)
	assertIsDir(err)
	_, err = d.WriteAt(nil, 5)
	assertIsDir(err)
	_, err = d.WriteString("")
	assertIsDir(err)
}

func TestDirReaddir(t *testing.T) {
	t.Parallel()
	i1 := memfs.FileInfo{
		Name: "bar",
	}
	i2 := memfs.FileInfo{
		Name: "baz",
	}
	i3 := memfs.FileInfo{
		Name: "boo",
	}
	d := memfs.NewDir("foo", dMode, dTime, []os.FileInfo{
		memfs.NewFileInfo(i1),
		memfs.NewFileInfo(i2),
		memfs.NewFileInfo(i3),
	})
	all, err := d.Readdir(0)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 3 {
		t.Fatal("was expecting 3")
	}
	some, err := d.Readdir(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(some) != 1 {
		t.Fatal("was expecting 1")
	}
	if actual := some[0].Name(); actual != i1.Name {
		t.Fatal("was expecting %s but got %s", i1.Name, actual)
	}
	some, err = d.Readdir(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(some) != 1 {
		t.Fatalf("was expecting 1, got %v", some)
	}
	if actual := some[0].Name(); actual != i2.Name {
		t.Fatalf("was expecting %s but got %s", i2.Name, actual)
	}
	some, err = d.Readdir(1)
	if err != io.EOF {
		t.Fatalf("was expecting io.EOF got %s", err)
	}
	if len(some) != 1 {
		t.Fatalf("was expecting 1, got %v", some)
	}
	if actual := some[0].Name(); actual != i3.Name {
		t.Fatalf("was expecting %s but got %s", i3.Name, actual)
	}
}

func TestDirReaddirnames(t *testing.T) {
	t.Parallel()
	i1 := memfs.FileInfo{
		Name: "bar",
	}
	i2 := memfs.FileInfo{
		Name: "baz",
	}
	i3 := memfs.FileInfo{
		Name: "boo",
	}
	d := memfs.NewDir("foo", dMode, dTime, []os.FileInfo{
		memfs.NewFileInfo(i1),
		memfs.NewFileInfo(i2),
		memfs.NewFileInfo(i3),
	})
	all, err := d.Readdirnames(0)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 3 {
		t.Fatal("was expecting 3")
	}
	some, err := d.Readdirnames(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(some) != 1 {
		t.Fatal("was expecting 1")
	}
	if actual := some[0]; actual != i1.Name {
		t.Fatal("was expecting %s but got %s", i1.Name, actual)
	}
	some, err = d.Readdirnames(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(some) != 1 {
		t.Fatalf("was expecting 1, got %v", some)
	}
	if actual := some[0]; actual != i2.Name {
		t.Fatalf("was expecting %s but got %s", i2.Name, actual)
	}
	some, err = d.Readdirnames(1)
	if err != io.EOF {
		t.Fatalf("was expecting io.EOF got %s", err)
	}
	if len(some) != 1 {
		t.Fatalf("was expecting 1, got %v", some)
	}
	if actual := some[0]; actual != i3.Name {
		t.Fatalf("was expecting %s but got %s", i3.Name, actual)
	}
}

func TestDirSetDirInfos(t *testing.T) {
	t.Parallel()
	i1 := memfs.FileInfo{
		Name: "bar",
	}
	i2 := memfs.FileInfo{
		Name: "baz",
	}
	d := memfs.NewDir("foo", dMode, dTime, []os.FileInfo{memfs.NewFileInfo(i1)})
	some, _ := d.Readdirnames(1)
	if actual := some[0]; actual != i1.Name {
		t.Fatal("was expecting %s but got %s", i1.Name, actual)
	}
	err := d.SetDirInfos([]os.FileInfo{memfs.NewFileInfo(i2)})
	if err != nil {
		t.Fatal(err)
	}
	some, _ = d.Readdirnames(1)
	if actual := some[0]; actual != i2.Name {
		t.Fatalf("was expecting %s but got %s", i2.Name, actual)
	}
}

func TestDirAddDirInfo(t *testing.T) {
	t.Parallel()
	i1 := memfs.FileInfo{
		Name: "bar",
	}
	i2 := memfs.FileInfo{
		Name: "baz",
	}
	d := memfs.NewDir("foo", dMode, dTime, []os.FileInfo{memfs.NewFileInfo(i1)})
	some, _ := d.Readdirnames(1)
	if actual := some[0]; actual != i1.Name {
		t.Fatal("was expecting %s but got %s", i1.Name, actual)
	}
	err := d.AddDirInfo(memfs.NewFileInfo(i2))
	if err != nil {
		t.Fatal(err)
	}
	some, _ = d.Readdirnames(2)
	if actual := some[1]; actual != i2.Name {
		t.Fatalf("was expecting %s but got %s", i2.Name, actual)
	}
}
