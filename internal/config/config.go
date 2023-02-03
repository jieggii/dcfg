package config

import (
	"encoding/json"
	"fmt"
	"github.com/jieggii/dcfg/internal/fs"
	"io"
	"os"
	"strings"
)

const configFilePermission = 0644

type Config struct {
	Bindings  Bindings  `json:"bindings"`  // path-to-path bindings
	Additions Additions `json:"additions"` // paths to additions
	Pinned    Pinned    `json:"pins"`      // paths to pinned objects
}

func (c *Config) ResolveAdditionDestination(addition string) (string, bool) {
	for _, bindingSource := range c.Bindings.Sources {
		if strings.HasPrefix(addition, bindingSource) {
			bindingDestination, err := c.Bindings.ResolveDestination(bindingSource)
			if err != nil {
				panic(err)
			}
			destination := strings.Replace(addition, bindingSource, bindingDestination, 1)
			return destination, true
		}
	}
	return "", false
}

func (c *Config) ResolveAdditionSource(additionDestination string) (string, bool) {
	for _, addition := range c.Additions.Paths {
		destination, resolved := c.ResolveAdditionDestination(addition)
		if resolved {
			if destination == additionDestination {
				return addition, true
			}
		}
	}
	return "", false
}

//func (c *Config) Validate() []error {
//	return nil
//}

func (c *Config) DumpToFile(path string) error {
	data, err := json.MarshalIndent(c, "", " ")

	if err != nil {
		panic(fmt.Errorf("unable to marshal config (%v)", err))
	}
	if err = os.WriteFile(path, data, configFilePermission); err != nil {
		return err
	}
	return nil
}

func NewConfig() *Config {
	return &Config{
		Bindings:  Bindings{},
		Additions: Additions{},
		Pinned:    Pinned{},
	}
}

func NewConfigFromFile(path string) (*Config, error) {
	exists, err := fs.NodeExists(path)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("dcfg config file '%v' does not exist", path)
	}

	cfg := NewConfig()

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open config file '%v' (%v)", path, err)
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}(file)

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read config file '%v' (%v)", path, err)
	}

	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("could not parse config file '%v' (%v)", path, err)
	}
	return cfg, nil
}
