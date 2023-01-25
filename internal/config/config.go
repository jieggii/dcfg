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

func (b Bindings) MarshalYAML() (any, error) {
	var bindingBatches []map[string]string

	if len(b.Sources) != len(b.Destinations) {
		panic("count of sources does not match count of destinations")
	}
	for i := range b.Sources {
		batch := make(map[string]string)
		batch[b.Sources[i]] = b.Destinations[i]
		bindingBatches = append(bindingBatches, batch)
	}
	return bindingBatches, nil
}

func (b *Bindings) AppendBinding(source string, destination string) error {
	if itemIsInArray(b.Sources, source) {
		return fmt.Errorf("path '%v' has already been registered as binding source", source)
	}
	b.Sources = append(b.Sources, source)

	if itemIsInArray(b.Destinations, destination) {
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

func (a Additions) MarshalYAML() (any, error) {
	return &a.Paths, nil
}

func (a *Additions) AppendAddition(path string) error {
	if itemIsInArray(a.Paths, path) {
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

func (p Pins) MarshalYAML() (any, error) {
	return &p.Paths, nil
}

func (p *Pins) AppendPin(path string) error {
	if itemIsInArray(p.Paths, path) {
		return fmt.Errorf("path '%v' is already pinned", path)
	}
	p.Paths = append(p.Paths, path)

	return nil
}

type Config struct {
	Context string `yaml:"context"` // path to context directory

	Bindings  Bindings  `yaml:"bindings,flow"`  // path-to-path bindings
	Additions Additions `yaml:"additions,flow"` // paths to additions
	Pins      Pins      `yaml:"pins,flow"`      // paths to pins
}

func (c *Config) DumpToFile(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		panic(fmt.Errorf("unable to marshal config (%v)", err))
	}
	if err = os.WriteFile(path, data, 0644); err != nil {
		return err
	}
	return nil
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

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("could not parse config %v (%v)", path, err)
	}
	return cfg, nil
}
