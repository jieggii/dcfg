package config

import (
	"fmt"
	"github.com/jieggii/dcfg/internal/util"
)

// Pinned is array of pinned paths.
// "pinned path" means path to file or directory which
// will not be removed when running `dcfg clean` command.
type Pinned []string

// Append appends `path` to array if pinned paths.
func (p *Pinned) Append(path string) error {
	if p.Exists(path) {
		return fmt.Errorf("'%v' is already pinned", path)
	}
	*p = append(*p, path)
	return nil
}

// Remove removes `path` from array os pinned paths.
func (p *Pinned) Remove(path string) error {
	i := util.ItemIndex(*p, path)
	if i == -1 {
		return fmt.Errorf("'%v' is not pinned", path)
	}
	*p = util.RemoveItem(*p, i)
	return nil
}

// Exists returns true if `path` is in array of pinned paths.
func (p *Pinned) Exists(path string) bool {
	return util.ItemIsInArray(*p, path)
}

// Any returns true if there is at least one path in array of pinned paths.
func (p *Pinned) Any() bool {
	return len(*p) != 0
}
