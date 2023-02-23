package commands

import (
	"fmt"
	"github.com/jieggii/dcfg/internal/config"
	"github.com/jieggii/dcfg/internal/fs"
	"github.com/jieggii/dcfg/internal/output"
	"github.com/urfave/cli/v2"
)

// Init command creates dcfg config file.
func Init(ctx *cli.Context) error {
	cfgPath := ctx.String("config")
	cfgExists, err := fs.NodeExists(cfgPath)
	if err != nil {
		return fmt.Errorf("could not check if dcfg config file '%v' exists (%v)",
			cfgPath, err,
		)
	}
	if cfgExists {
		return fmt.Errorf("dcfg config file '%v' already exists", cfgPath)
	}
	cfg := config.NewConfig()
	if err := cfg.DumpToFile(cfgPath); err != nil {
		return err
	}
	output.Plus.Printf("created dcfg config file '%v'", cfgPath)
	return nil
}
