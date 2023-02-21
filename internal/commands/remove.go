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
		var source string      // path to source of target
		var destination string // path to destination of collected target

		target = path.Clean(target)

		if path.IsAbs(target) { // target source is provided
			source = target
			destination, _ = cfg.ResolveTargetDestination(target)
			// don't check if target destination path is resolved because
			// suitable binding can be already deleted and resolving will
			// fail because of that
		} else { // path to target destination is provided
			var resolved bool
			source, resolved = cfg.ResolveTargetSource(target)
			if !resolved {
				return fmt.Errorf("could not resolve source of '%v'", target)
			}

			destination = target
		}
		if err := cfg.Targets.Remove(target); err != nil {
			output.Error.Printf(
				"could not remove '%v' from targets list",
				source,
			)
			continue
		}
		output.Minus.Printf("removed '%v' from targets list", target)

		targetDestinationExists, err := fs.NodeExists(destination)
		if err != nil {
			output.Warning.Println(
				"could not check if '%v' exists (%v)", destination, err,
			)
		}

		if !soft && targetDestinationExists {
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
