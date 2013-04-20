package pkgfs

import (
	"testing"
)

// /file
// /dir/file
// /dir : /dir/file1 /dir/file2
// missing prefix slash
// ../file
// / with glob
// empty dir
// subdir with glob with real files but not matching the glob
// empty string
// non file/dir (like socket/pipe)
//

func TestStub(t *testing.T) {
	t.Parallel()
}
