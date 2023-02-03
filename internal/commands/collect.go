package commands

import (
	"github.com/jieggii/dcfg/internal/config"
	"github.com/jieggii/dcfg/internal/fs"
	"github.com/jieggii/dcfg/internal/output"
	"github.com/urfave/cli/v2"
	"path"
)

func Collect(ctx *cli.Context) error {
	cfgPath := ctx.String("config")

	cfg, err := config.NewConfigFromFile(cfgPath)
	if err != nil {
		return err
	}

	for _, addition := range cfg.Additions.Paths {
		destination, resolved := cfg.ResolveAdditionDestination(addition)
		if !resolved {
			output.Warning.Printf("not resolved %v\n", destination)
			continue
		}
		destination = path.Join(cfg.Context.String(), destination)
		if err := fs.Copy(addition, destination); err != nil {
			output.Warning.Printf("could not copy %v to %v (%v)\n", addition, destination, err)
			continue
		}
		output.Plus.Printf("copied addition %v to %v\n", addition, destination)
	}
	return nil
}
