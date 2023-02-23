package fs

import (
	cp "github.com/otiai10/copy"
	"os"
)

var copyOptions = cp.Options{
	OnSymlink: func(_ string) cp.SymlinkAction {
		return cp.Deep // create hard copy of contents when reach symlink
	},
	OnDirExists: func(_, _ string) cp.DirExistsAction {
		return cp.Replace // replace content of directory if it already exists
	},
}

// Copy copies `source` to `destination`, no matter is it directory
// containing other directories or file.
func Copy(source string, destination string) error {
	return cp.Copy(source, destination, copyOptions)
}

// NodeExists returns true if file or directory at `path` exists.
func NodeExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
