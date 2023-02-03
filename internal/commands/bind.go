package commands

import (
	"fmt"
	"github.com/jieggii/dcfg/internal/config"
	"github.com/jieggii/dcfg/internal/output"
	"github.com/urfave/cli/v2"
	"path"
)

func Bind(ctx *cli.Context) error {
	cfgPath := ctx.String("config")
	remove := ctx.Bool("remove")

	cfg, err := config.NewConfigFromFile(cfgPath)
	if err != nil {
		return err
	}
	args := ctx.Args()
	argsCount := args.Len()

	source := path.Clean(args.Get(0))
	if !path.IsAbs(source) {
		return fmt.Errorf("source path must be absolute (got relative path '%v')", source)
	}

	if remove {
		if argsCount != 1 {
			return fmt.Errorf("bind command with '--remove' flag takes only 1 argument (got %v) - source.\nusage: %v", argsCount, ctx.Command.UsageText)
		}
		destination, err := cfg.Bindings.ResolveDestination(source)
		if err != nil {
			return err
		}
		if err := cfg.Bindings.Remove(source); err != nil {
			return err
		}
		if err = cfg.DumpToFile(cfgPath); err != nil {
			return err
		}
		output.Minus.Printf("removed binding: %v -> %v.\n", source, destination)

	} else {
		if argsCount != 2 {
			return fmt.Errorf("bind command without '--remove' flag takes exactly 2 arguments (got %v).\nusage: %v", argsCount, ctx.Command.UsageText)
		}
		destination := path.Clean(args.Get(1))
		if path.IsAbs(destination) {
			return fmt.Errorf("destination path must be relative (got absolute path '%v')", destination)
		}
		if err = cfg.Bindings.Append(source, destination); err != nil {
			return err
		}
		if err = cfg.DumpToFile(cfgPath); err != nil {
			return err
		}
		output.Plus.Printf("registered new binding: %v -> %v.\n", source, destination)
	}

	return nil
}
