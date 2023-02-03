package fs

import (
	cp "github.com/otiai10/copy"
	"os"
)

func Copy(addition string, destination string) error {
	options := cp.Options{
		OnSymlink: func(_ string) cp.SymlinkAction {
			return cp.Deep
		},
		OnDirExists: func(_, _ string) cp.DirExistsAction {
			return cp.Replace
		},
	}
	return cp.Copy(addition, destination, options)
}

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
