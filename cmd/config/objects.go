package config

import (
	p "path"
	"strings"
)

type Bindings struct {
	Sources      []string
	Destinations []string

	LongestSourceLength int // needed for pretty output
}

func (bindings *Bindings) appendBinding(source string, destination string) {
	sourceLen := len(source)
	if sourceLen > bindings.LongestSourceLength {
		bindings.LongestSourceLength = sourceLen
	}
	bindings.Sources = append(bindings.Sources, source)
	bindings.Destinations = append(bindings.Destinations, destination)
}

func newBindings() Bindings {
	return Bindings{LongestSourceLength: 0}
}

type Additions struct {
	Paths []string // list of additions

	LongestPathLength int // needed for pretty output

}

func (additions *Additions) appendAddition(path string) {
	pathLen := len(path)
	if pathLen > additions.LongestPathLength {
		additions.LongestPathLength = pathLen
	}
	additions.Paths = append(additions.Paths, path)
}

func newAdditions() Additions {
	return Additions{LongestPathLength: 0}
}

type Config struct {
	Bindings  Bindings
	Additions Additions
	Pins      []string

	Context       string
	ContextWasSet bool
}

func (config *Config) AppendPin(path string) {
	config.Pins = append(config.Pins, path)
}

func (config *Config) PathIsDestination(path string) bool {
	for _, destination := range config.Bindings.Destinations {
		if strings.HasSuffix(p.Clean(destination), p.Clean(path)) { // todo test
			return true
		}
	}
	return false
}

func (config *Config) PathIsPinned(path string) bool {
	for _, pin := range config.Pins {
		if strings.HasSuffix(p.Clean(pin), p.Clean(path)) { // todo test
			return true
		}
	}
	return false
}

func newConfig() *Config {
	return &Config{
		Bindings:      newBindings(),
		Additions:     newAdditions(),
		Context:       "./",
		ContextWasSet: false,
	}
}
