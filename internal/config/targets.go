package config

import (
	"encoding/json"
	"fmt"
	"github.com/jieggii/dcfg/internal/fs"
	"github.com/jieggii/dcfg/internal/util"
)

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

func (a *Targets) Append(path string) error {
	if a.Exists(path) {
		return fmt.Errorf("path '%v' is already registered as addition", path)
	}
	a.Paths = append(a.Paths, path)

	pathLen := len(path)
	if pathLen > a.LongestPathLen {
		a.LongestPathLen = pathLen
	}

	return nil
}

func (a *Targets) Remove(path string) error {
	i := util.ItemIndex(a.Paths, path)
	if i == -1 {
		return fmt.Errorf("'%v' is not registered as addition", path)
	}
	a.Paths = util.RemoveItem(a.Paths, i)
	return nil
}

func (a *Targets) Exists(path string) bool {
	return util.ItemIsInArray(a.Paths, path)
}

func (a *Targets) Any() bool {
	return len(a.Paths) != 0
}

func (a *Targets) IsCollected(destination string) (bool, error) {
	collected, err := fs.NodeExists(destination)
	if err != nil {
		return false, fmt.Errorf("could not check if '%v' exists (%v)", destination, err)
	}
	return collected, nil
}
