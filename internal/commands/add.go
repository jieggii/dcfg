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

	cfg, err := config.NewConfigFromFile(cfgPath)
	if err != nil {
		return err
	}
	addition := path.Clean(ctx.Args().First())
	if !path.IsAbs(addition) {
		return fmt.Errorf("path to addition must be absolute (got relative path '%v')", addition)
	}
	if err := cfg.Additions.Append(addition); err != nil {
		return err
	}
	if err := cfg.DumpToFile(cfgPath); err != nil {
		return err
	}
	output.Plus.Printf("appended new addition %v.", addition)

	if collect {
		destination, resolved := cfg.ResolveAdditionDestination(addition)
		if !resolved {
			return errors.New("didn't collect because of missing suitable binding for the addition")
		}
		destination = path.Join(cfg.Context.String(), destination)
		if err := fs.Copy(addition, destination); err != nil {
			return fmt.Errorf("could not copy addition to '%v' (%v)", destination, err)
		}
		output.Plus.Printf("copied '%v' to '%v' (in '%v')", destination, cfg.Context)
	}
	return nil
}
