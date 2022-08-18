package config

type Bindings struct {
	Sources      []string
	Destinations []string

	DestinationPrefix       string // prefix which will be joined to destination path
	DestinationPrefixWasSet bool   // indicates if d.p. was set in the config file
	LongestSourceLength     int    // needed for pretty output
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
	return Bindings{LongestSourceLength: 0, DestinationPrefix: "./"}
}

type Additions struct {
	LongestPathLength int      // needed for pretty output
	Paths             []string // list of additions
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
}

func newConfig() *Config {
	return &Config{
		Bindings:  newBindings(),
		Additions: newAdditions(),
	}
}
