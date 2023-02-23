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

// Remove command removes target.
func Remove(ctx *cli.Context) error {
	cfgPath := ctx.String("config")
	cfg, err := config.NewConfigFromFile(cfgPath)
	if err != nil {
		return err
	}

	// flags
	soft := ctx.Bool("soft")

	paths := ctx.Args().Slice()

	for _, anyPath := range paths {
		var source string      // path to source of target
		var destination string // path to destination of collected target

		anyPath = path.Clean(anyPath)

		if path.IsAbs(anyPath) { // target source is provided
			source = anyPath
			destination, _ = cfg.ResolveTargetDestination(anyPath)
			// don't check if target destination path is resolved because
			// suitable binding can be already deleted and resolving will
			// fail because of that
		} else { // path to target destination is provided
			var resolved bool
			source, resolved = cfg.ResolveTargetSource(anyPath)
			if !resolved {
				return fmt.Errorf("could not resolve source of '%v'", anyPath)
			}

			destination = anyPath
		}
		if err := cfg.Targets.Remove(source); err != nil {
			output.Warning.Println(err)
			continue
		}
		output.Minus.Printf("removed '%v' from targets list", source)

		targetDestinationExists, err := fs.NodeExists(destination)
		if err != nil {
			output.Warning.Println(
				"could not check if '%v' exists (%v)", destination, err,
			)
		}

		if !soft && targetDestinationExists {
			err := os.RemoveAll(destination)
			if err != nil {
				output.Warning.Printf(
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
