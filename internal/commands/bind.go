package commands

import (
	"github.com/jieggii/dcfg/internal/config"
	"github.com/jieggii/dcfg/internal/output"
	"github.com/urfave/cli/v2"
)

func Bind(ctx *cli.Context) error {
	cfgPath := ctx.String("config")
	cfg, err := config.NewConfigFromFile(cfgPath)
	if err != nil {
		return err
	}
	args := ctx.Args()
	source := args.Get(0)
	destination := args.Get(1)
	// todo: validate src and dest
	// ...
	if err = cfg.Bindings.AppendBinding(source, destination); err != nil {
		return err
	}
	if err = cfg.DumpToFile(cfgPath); err != nil {
		return err
	}
	output.Success.Printf("Added new binding: %v -> %v.\n", source, destination)
	return nil
}
