package commands

import (
	"fmt"
	"github.com/jieggii/dcfg/internal/config"
	"github.com/jieggii/dcfg/internal/fs"
	"github.com/jieggii/dcfg/internal/output"
	"github.com/urfave/cli/v2"
	"os"
	"path"
)

func Remove(ctx *cli.Context) error {
	cfgPath := ctx.String("config")
	cfg, err := config.NewConfigFromFile(cfgPath)
	if err != nil {
		return err
	}

	// flags
	soft := ctx.Bool("soft")

	targets := ctx.Args().Slice()
	for _, target := range targets {
		var addition string    // addition = path to source
		var destination string // path to destination of collected addition

		target = path.Clean(target)

		if path.IsAbs(target) { // addition is provided
			addition = target
			destination, _ = cfg.ResolveAdditionDestination(target)
			// don't check if addition destination path is resolved because
			// suitable binding can be already deleted and resolving will
			// fail because of that
		} else { // path to addition destination is provided
			var resolved bool
			addition, resolved = cfg.ResolveAdditionSource(target)
			if !resolved {
				return fmt.Errorf("could not resolve source of '%v'", target)
			}
			destination = target
		}
		if err := cfg.Additions.Remove(addition); err != nil {
			output.Error.Printf(
				"could not remove '%v' from additions list",
				addition,
			)
			continue
		}
		output.Minus.Printf("removed '%v' from additions list", addition)

		additionDestinationExists, err := fs.NodeExists(destination)
		if err != nil {
			output.Warning.Println(
				"could not check if '%v' exists (%v)", destination, err,
			)
		}

		if !soft && additionDestinationExists {
			err := os.RemoveAll(destination)
			if err != nil {
				output.Error.Printf(
					"could not remove '%v' (%v)", destination, err,
				)
				continue
			}
			output.Minus.Printf("removed '%v'", destination)
		}
	}
	if err := cfg.DumpToFile(cfgPath); err != nil {
		return err
	}

	return nil
}
