package commands

import (
	"fmt"
	"github.com/jieggii/dcfg/internal/config"
	"github.com/jieggii/dcfg/internal/output"
	"github.com/urfave/cli/v2"
	"os"
	"path"
)

func Remove(ctx *cli.Context) error {
	cfgPath := ctx.String("config")
	soft := ctx.Bool("soft")

	cfg, err := config.NewConfigFromFile(cfgPath)
	if err != nil {
		return err
	}

	var additionSourcePath string
	var additionDestinationPath string
	addition := path.Clean(ctx.Args().First())

	if path.IsAbs(addition) { // addition source is provided
		additionSourcePath = addition
		additionDestinationPath, _ = cfg.ResolveAdditionDestination(addition)
		// don't check if addition destination path is resolved 'cause we
		// don't need it
	} else { // path to destination is provided
		var resolved bool
		additionSourcePath, resolved = cfg.ResolveAdditionSource(addition)
		if !resolved {
			return fmt.Errorf("could not resolve source of '%v'", addition)
		}
		additionDestinationPath = addition
	}
	if err := cfg.Additions.Remove(additionSourcePath); err != nil {
		return err
	}
	if err := cfg.DumpToFile(cfgPath); err != nil {
		return err
	}
	output.Minus.Printf("removed addition '%v' from additions list.", additionSourcePath)
	if !soft && additionDestinationPath != "" {
		err := os.RemoveAll(additionDestinationPath)
		if err != nil {
			return err
		}
		output.Minus.Printf("removed '%v'.", additionDestinationPath)
	}
	return nil
}
