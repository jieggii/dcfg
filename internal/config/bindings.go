package config

import (
	"encoding/json"
	"fmt"
	"github.com/jieggii/dcfg/internal/util"
	"strings"
)

const bindingWithSourceDoesNotExistErrorText = "binding with source '%v' does not exist"

// Bindings represents bindings (yes).
//
// Binding is like mount point or alias, they map source path
// to destination path.
//
// For example, binding `/home/user -> ./user-home` means that all targets
// from `/home/user` will be copied to `./user-home`
type Bindings struct {
	// Using two arrays instead of map[string]string because maps in Go
	// don't store keys order (see https://go.dev/blog/maps#iteration-order).
	Sources      []string // list of absolute paths to sources
	Destinations []string // list of relative paths to destinations

	LongestSourceLen int // length of the longest source path, needed for pretty-printing
}

func (b Bindings) MarshalJSON() ([]byte, error) {
	var bindingBatches []map[string]string

	if len(b.Sources) != len(b.Destinations) {
		panic("count of sources does not match count of destinations")
	}

	for i := range b.Sources {
		batch := make(map[string]string)
		batch[b.Sources[i]] = b.Destinations[i]
		bindingBatches = append(bindingBatches, batch)
	}
	data, err := json.Marshal(bindingBatches)
	return data, err
}

func (b *Bindings) UnmarshalJSON(data []byte) error {
	var bindingBatches []map[string]string
	if err := json.Unmarshal(data, &bindingBatches); err != nil {
		return err
	}
	for i := range bindingBatches {
		for source, destination := range bindingBatches[i] {
			err := b.Append(source, destination)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Append registers new binding.
func (b *Bindings) Append(source string, destination string) error {
	if b.SourceExists(source) {
		return fmt.Errorf("'%v' is already registered as binding source", source)
	}
	b.Sources = append(b.Sources, source)

	sourceLength := len(source)
	if sourceLength > b.LongestSourceLen {
		b.LongestSourceLen = sourceLength
	}

	// not checking if destination exists because same destination may be used with multiple sources
	b.Destinations = append(b.Destinations, destination)

	return nil
}

// Remove removes binding with source `source`.
func (b *Bindings) Remove(source string) error {
	i := util.ItemIndex(b.Sources, source)
	if i == -1 {
		return fmt.Errorf(bindingWithSourceDoesNotExistErrorText, source)
	}
	b.Sources = util.RemoveItem(b.Sources, i)
	b.Destinations = util.RemoveItem(b.Destinations, i)
	return nil
}

// ResolveDestination returns destination of binding with source `source`.
func (b *Bindings) ResolveDestination(source string) (string, error) {
	for i := range b.Sources {
		if b.Sources[i] == source {
			return b.Destinations[i], nil
		}
	}
	return "", fmt.Errorf(bindingWithSourceDoesNotExistErrorText, source)
}

// SourceExists returns true if `path` is registered as a binging source.
func (b *Bindings) SourceExists(path string) bool {
	return util.ItemIsInArray(b.Sources, path)
}

// DestinationExists returns true if `path` is registered as a binging destination.
func (b *Bindings) DestinationExists(path string) bool {
	return util.ItemIsInArray(b.Destinations, path)
}

// DestinationWithPrefixExists returns true if path with prefix `prefix` is registered
// as a binding destination.
func (b *Bindings) DestinationWithPrefixExists(prefix string) bool {
	for _, destination := range b.Destinations {
		if strings.HasPrefix(destination, prefix) {
			return true
		}
	}
	return false
}

// Any returns true if there is at least one binding registered.
func (b *Bindings) Any() bool {
	return len(b.Sources) != 0
}
