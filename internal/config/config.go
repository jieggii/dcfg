package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const configFilePermission = 0644

type Config struct {
	Context Context `json:"context"` // path to context directory

	Bindings  Bindings  `json:"bindings"`  // path-to-path bindings
	Additions Additions `json:"additions"` // paths to additions
	Pinned    Pinned    `json:"pins"`      // paths to pinned objects
}

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
		Context:   ".",
		Bindings:  Bindings{},
		Additions: Additions{},
		Pinned:    Pinned{},
	}
}

func NewConfigFromFile(path string) (*Config, error) {
	cfg := NewConfig()

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open config %v (%v)", path, err)
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}(file)

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read config %v (%v)", path, err)
	}

	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("could not parse config %v (%v)", path, err)
	}
	return cfg, nil
}
