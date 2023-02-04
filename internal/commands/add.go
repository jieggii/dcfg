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

func Add(ctx *cli.Context) error {
	cfgPath := ctx.String("config")

	collect := ctx.Bool("collect")
	verbose := ctx.Bool("verbose")

	cfg, err := config.NewConfigFromFile(cfgPath)
	if err != nil {
		return err
	}

	addition := path.Clean(ctx.Args().First())
	if !path.IsAbs(addition) {
		return fmt.Errorf("path to addition must be absolute (got relative path '%v')", addition)
	}

	destination, resolved := cfg.ResolveAdditionDestination(addition)
	if !resolved {
		return errors.New("could not resolve addition destination (missing suitable binding)")
	}
	if err := cfg.Additions.Append(addition); err != nil {
		return err
	}
	if err := cfg.DumpToFile(cfgPath); err != nil {
		return err
	}
	output.Plus.Printf("appended new addition '%v'\n", addition)
	if verbose {
		output.Verbose.Printf("the addition will be copied to '%v' when collecting\n", destination)
	}

	if collect {
		if err := fs.Copy(addition, destination); err != nil {
			return fmt.Errorf("could not copy '%v' to '%v' (%v)", addition, destination, err)
		}
		output.Plus.Printf("copied '%v' -> '%v'", addition, destination)
	}
	return nil
}
