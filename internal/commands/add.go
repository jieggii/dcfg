package commands

import (
	"errors"
	"fmt"
	"github.com/jieggii/dcfg/internal/config"
	"github.com/jieggii/dcfg/internal/fs"
	"github.com/jieggii/dcfg/internal/output"
	"github.com/urfave/cli/v2"
	"path"
)

// Add command adds new target.
func Add(ctx *cli.Context) error {
	cfgPath := ctx.String("config")
	cfg, err := config.NewConfigFromFile(cfgPath)
	if err != nil {
		return err
	}

	// options
	collect := ctx.Bool("collect")

	// arguments
	targets := ctx.Args().Slice()

	destinations := map[string]string{}

	// validating:
	for _, target := range targets {
		target = path.Clean(target)
		if !path.IsAbs(target) {
			return fmt.Errorf(
				"path to target must be absolute (got relative path '%v')",
				target,
			)
		}
		exists, err := fs.NodeExists(target)
		if err != nil {
			return err
		}
		if !exists {
			return errors.New("'%v' does not exist")
		}

		destination, resolved := cfg.ResolveTargetDestination(target)
		if !resolved {
			return errors.New(
				"could not resolve target destination (missing suitable binding)",
			)
		}
		destinations[target] = destination
	}

	// performing actions:
	for _, target := range targets {
		target = path.Clean(target)
		destination := destinations[target]

		if err := cfg.Targets.Append(target); err != nil {
			output.Error.Println(err)
			continue
		}
		output.Plus.Printf("appended new target '%v'", target)

		if collect {
			if err := fs.Copy(target, destination); err != nil {
				output.Error.Printf(
					"could not copy '%v' to '%v' (%v)",
					target, destination, err,
				)
				continue
			}
			output.Plus.Printf(
				"copied '%v' -> '%v'", target, destination,
			)
		}
	}

	if err := cfg.DumpToFile(cfgPath); err != nil {
		return err
	}

	return nil
}
