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
	Bindings Bindings `json:"bindings"` // path-to-path bindings
	Targets  Targets  `json:"targets"`  // paths to destinations of targets
	Pinned   Pinned   `json:"pins"`     // paths to pinned objects
}

func (c *Config) ResolveTargetDestination(target string) (string, bool) {
	for _, bindingSource := range c.Bindings.Sources {
		if strings.HasPrefix(target, bindingSource) {
			bindingDestination, err := c.Bindings.ResolveDestination(bindingSource)
			if err != nil {
				panic(err)
			}
			destination := strings.Replace(target, bindingSource, bindingDestination, 1)
			return destination, true
		}
	}
	return "", false
}

func (c *Config) ResolveTargetSource(targetDestination string) (string, bool) {
	for _, target := range c.Targets.Paths {
		destination, resolved := c.ResolveTargetDestination(target)
		if resolved {
			if destination == targetDestination {
				return target, true
			}
		}
	}
	return "", false
}

func (c *Config) DumpToFile(path string) error {
	data, err := json.MarshalIndent(c, "", "    ")
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
		Bindings: Bindings{},
		Targets:  Targets{},
		Pinned:   Pinned{},
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
