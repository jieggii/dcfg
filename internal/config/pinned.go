package config

import "fmt"

type Pinned []string

func (p *Pinned) Append(path string) error {
	if p.Exists(path) {
		return fmt.Errorf("path '%v' is already pinned", path)
	}
	*p = append(*p, path)
	return nil
}

func (p *Pinned) Remove(path string) error {
	i := itemIndex(*p, path)
	if i == -1 {
		return fmt.Errorf("path '%v' is not pinned", path)
	}
	*p = removeItem(*p, i)
	return nil
}

func (p *Pinned) Exists(path string) bool {
	return itemIsInArray(*p, path)
}

func (p *Pinned) IsPresent() bool {
	return len(*p) != 0
}
