package diff

import (
	"os/exec"
	"syscall"
)

var diffFlags = []string{
	"--unified",
	"--recursive",
	"--color=always",
	//"--report-identical-files",
}

// GetDiff returns combined output (stdout + stderr) of diff run over `file1` and `file2`.
func GetDiff(diffBinPath string, file1 string, file2 string) (string, error) {
	flags := diffFlags
	flags = append(append(flags, file1), file2)
	output, err := exec.Command(diffBinPath, flags...).CombinedOutput()
	if err != nil {
		switch e := err.(type) {
		case *exec.ExitError: // run diff and it exited
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
