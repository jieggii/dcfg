package fs

import (
	"crypto/sha256"
	cp "github.com/otiai10/copy"
	"io"
	"os"
)

var copyOptions = cp.Options{
	OnSymlink: func(_ string) cp.SymlinkAction {
		return cp.Deep
	},
	OnDirExists: func(_, _ string) cp.DirExistsAction {
		return cp.Replace
	},
}

func Copy(source string, destination string) error {
	return cp.Copy(source, destination, copyOptions)
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

func CalculateFileHash(path string) (error, []byte) {
	file, err := os.Open(path)
	if err != nil {
		return err, nil
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}(file)
	h := sha256.New224()
	if _, err := io.Copy(h, file); err != nil {
		return err, nil
	}
	return nil, h.Sum(nil)
}

// file hash consists of file content
// dir hash consists of path to dir
// tree hash consists of

//func CalculateNodeHash(path string) (error, []byte) {
//	info, err := os.Stat(path)
//	if err != nil {
//		return err, nil
//	}
//	if !info.IsDir() { // `path` is a file
//		return calculateFileHash(path)
//	} else {
//		var paths []string
//		err := filepath.WalkDir(path, func(curPath string, d fs.DirEntry, err error) error {
//			if err != nil {
//				return err
//			}
//			paths = append(paths, curPath)
//			return nil
//		})
//		if err != nil {
//			return err, nil
//		}
//
//	}
//	return nil, nil
//}
