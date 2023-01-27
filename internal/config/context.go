package config

import (
	"encoding/json"
	"fmt"
	"path"
)

type Context string

func (c *Context) String() string {
	return string(*c)
}

// UnmarshalJSON simply cleans context path (using path.Clean) before storing it
func (c *Context) UnmarshalJSON(data []byte) error {
	var context string

	if err := json.Unmarshal(data, &context); err != nil {
		return err
	}
	context = path.Clean(context)
	if path.IsAbs(context) {
		return fmt.Errorf("context path must be relative (got absolute '%v')", context)
	}
	*c = Context(context)
	return nil
}
