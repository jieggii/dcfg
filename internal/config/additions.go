package config

import "fmt"

type Additions []string

func (a *Additions) Append(path string) error {
	if a.Exists(path) {
		return fmt.Errorf("path '%v' is already registered as addition", path)
	}
	*a = append(*a, path)
	return nil
}

func (a *Additions) Remove(path string) error {
	i := itemIndex(*a, path)
	if i == -1 {
		return fmt.Errorf("path '%v' is not registered as addition", path)
	}
	*a = removeItem(*a, i)
	return nil
}

func (a *Additions) Exists(path string) bool {
	return itemIsInArray(*a, path)
}
