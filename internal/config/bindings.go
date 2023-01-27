package config

import (
	"encoding/json"
	"fmt"
)

const bindingWithSourceDoesNotExistErrorText = "binding with source '%v' does not exist"

type Bindings struct {
	Sources      []string // list of absolute paths to sources
	Destinations []string // list of relative to Config.Context paths to destinations
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
	// don't use simple map[string]string in this case because map
	// doesn't store keys order (see https://go.dev/blog/maps#iteration-order)
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

func (b *Bindings) Append(source string, destination string) error {
	if b.SourceExists(source) {
		return fmt.Errorf("path '%v' is already registered as binding source", source)
	}
	b.Sources = append(b.Sources, source)

	// not checking if destination exists 'cause same destination may be used with multiple sources
	b.Destinations = append(b.Destinations, destination)
	return nil
}

func (b *Bindings) Remove(source string) error {
	i := itemIndex(b.Sources, source)
	if i == -1 {
		return fmt.Errorf(bindingWithSourceDoesNotExistErrorText, source)
	}
	b.Sources = removeItem(b.Sources, i)
	b.Destinations = removeItem(b.Destinations, i)
	return nil
}

func (b *Bindings) ResolveDestination(source string) (string, error) {
	for i := range b.Sources {
		if b.Sources[i] == source {
			return b.Destinations[i], nil
		}
	}
	return "", fmt.Errorf(bindingWithSourceDoesNotExistErrorText, source)
}

func (b *Bindings) SourceExists(path string) bool {
	return itemIsInArray(b.Sources, path)
}

func (b *Bindings) DestinationExists(path string) bool {
	return itemIsInArray(b.Destinations, path)
}
