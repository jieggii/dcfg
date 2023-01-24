package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

const defaultContext = "./"

type Bindings struct {
	Sources      []string // list of global paths to sources
	Destinations []string // list of relative to Config.Context paths to destinations
}

func (b *Bindings) UnmarshalYAML(unmarshal func(any) error) error {
	// don't use simple map[string]string in this case because map
	// doesn't store keys order (see https://go.dev/blog/maps#iteration-order)
	var bindingBatches []map[string]string
	if err := unmarshal(&bindingBatches); err != nil {
		return err
	}

	for i := range bindingBatches {
		for source, destination := range bindingBatches[i] {
			err := b.AppendBinding(source, destination)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Bindings) AppendBinding(source string, destination string) error {
	if itemExists(b.Sources, source) {
		return fmt.Errorf("path '%v' has already been registered as binding source", source)
	}
	b.Sources = append(b.Sources, source)

	if itemExists(b.Destinations, destination) {
		return fmt.Errorf("path '%v' has already been registered as binding destination", destination)
	}
	b.Destinations = append(b.Destinations, destination)

	return nil
}

type Additions struct {
	Paths []string // list of global paths to additions
}

func (a *Additions) UnmarshalYAML(unmarshal func(any) error) error {
	var additions []string
	if err := unmarshal(&additions); err != nil {
		return err
	}
	for _, addition := range additions {
		if err := a.AppendAddition(addition); err != nil {
			return err
		}
	}

	return nil
}

func (a *Additions) AppendAddition(path string) error {
	if itemExists(a.Paths, path) {
		return fmt.Errorf("path '%v' has already been added", path)
	}
	a.Paths = append(a.Paths, path)

	return nil
}

type Pins struct {
	Paths []string // list of relative to Config.Context paths to pins
}

func (p *Pins) UnmarshalYAML(unmarshal func(any) error) error {
	var pins []string
	if err := unmarshal(&pins); err != nil {
		return err
	}

	for _, pin := range pins {
		if err := p.AppendPin(pin); err != nil {
			return err
		}
	}

	return nil
}

func (p *Pins) AppendPin(path string) error {
	if itemExists(p.Paths, path) {
		return fmt.Errorf("path '%v' is already pinned", path)
	}
	p.Paths = append(p.Paths, path)

	return nil
}

type Config struct {
	Context string `yaml:"context"` // path to context directory

	Bindings  Bindings  `yaml:"bindings"`  // path-to-path bindings
	Additions Additions `yaml:"additions"` // paths to additions
	Pins      Pins      `yaml:"pins"`      // paths to pins
}

func NewConfig() *Config {
	return &Config{
		Context:   defaultContext,
		Bindings:  Bindings{},
		Additions: Additions{},
		Pins:      Pins{},
	}
}

func NewConfigFromFile(path string) (*Config, error) {
	cfg := NewConfig()

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening config %v: %v", path, err)
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}(file)

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading config %v: %v", path, err)
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("error parsing config %v: %v", path, err)
	}
	return cfg, nil
}
