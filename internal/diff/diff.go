package diff

import (
	"os/exec"
	"syscall"
)

const diffPath = "/usr/bin/diff"

var diffFlags = []string{
	"--report-identical-files",
	"--unified",
	"--recursive",
	"--color=always",
}

func GetDiff(file1 string, file2 string) (string, error) {
	flags := diffFlags
	flags = append(append(flags, file1), file2)
	output, err := exec.Command(diffPath, flags...).CombinedOutput()
	if err != nil {
		switch e := err.(type) {
		case *exec.ExitError: // could run diff and it exited
			if status, ok := e.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() != 1 { // diff returned error
					return string(output), err
				}
			}
		default: // could not run diff
			return "", err
		}
	}
	return string(output), nil
}
