package config

import (
	"encoding/json"
	"fmt"
	"github.com/jieggii/dcfg/internal/fs"
	"github.com/jieggii/dcfg/internal/util"
)

// Targets represents array of targets.
// Target is a file or directory to be collected using dcfg.
type Targets struct {
	Paths []string // absolute paths to targets

	LongestPathLen int
}

func (a Targets) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(a.Paths)
	return data, err
}

func (a *Targets) UnmarshalJSON(data []byte) error {
	var paths []string

	if err := json.Unmarshal(data, &paths); err != nil {
		return err
	}
	for _, path := range paths {
		err := a.Append(path)
		if err != nil {
			return err
		}
	}
	return nil
}

// Append appends `path` to array of targets.
func (a *Targets) Append(path string) error {
	if a.Exists(path) {
		return fmt.Errorf("'%v' is already registered as target", path)
	}
	a.Paths = append(a.Paths, path)

	pathLen := len(path)
	if pathLen > a.LongestPathLen {
		a.LongestPathLen = pathLen
	}

	return nil
}

// Remove removes `path` from array of targets.
func (a *Targets) Remove(path string) error {
	i := util.ItemIndex(a.Paths, path)
	if i == -1 {
		return fmt.Errorf("'%v' is not registered as target", path)
	}
	a.Paths = util.RemoveItem(a.Paths, i)
	return nil
}

// Exists returns true if `path` is present in array of targets.
func (a *Targets) Exists(path string) bool {
	return util.ItemIsInArray(a.Paths, path)
}

// Any returns true if there is at least one target in arraay of targets.
func (a *Targets) Any() bool {
	return len(a.Paths) != 0
}

// IsCollected returns true if `destination` exists.
func (a *Targets) IsCollected(destination string) (bool, error) {
	collected, err := fs.NodeExists(destination)
	if err != nil {
		return false, fmt.Errorf("could not check if '%v' exists (%v)", destination, err)
	}
	return collected, nil
}
