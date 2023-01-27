package fs

import "strings"

func PathHasUserHomeDirPrefix(path string) bool {
	return strings.HasPrefix(path, "/home/")
}
