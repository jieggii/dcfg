package config

import (
	"fmt"
	"github.com/jieggii/dcfg/internal/util"
)

type Pinned []string

func (p *Pinned) Append(path string) error {
	if p.Exists(path) {
		return fmt.Errorf("'%v' is already pinned", path)
	}
	*p = append(*p, path)
	return nil
}

func (p *Pinned) Remove(path string) error {
	i := util.ItemIndex(*p, path)
	if i == -1 {
		return fmt.Errorf("'%v' is not pinned", path)
	}
	*p = util.RemoveItem(*p, i)
	return nil
}

func (p *Pinned) Exists(path string) bool {
	return util.ItemIsInArray(*p, path)
}

func (p *Pinned) Any() bool {
	return len(*p) != 0
}
