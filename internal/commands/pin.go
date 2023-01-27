package commands

import (
	"fmt"
	"github.com/jieggii/dcfg/internal/config"
	"github.com/jieggii/dcfg/internal/output"
	"github.com/urfave/cli/v2"
	"path"
	"strings"
)

func Pin(ctx *cli.Context) error {
	cfgPath := ctx.String("config")
	remove := ctx.Bool("remove")

	cfg, err := config.NewConfigFromFile(cfgPath)
	if err != nil {
		return err
	}
	args := ctx.Args()
	pinned := path.Clean(args.Get(0))

	if path.IsAbs(pinned) {
		return fmt.Errorf("path to folder to pin must be relative (got absolute path '%v')", pinned)
	}

	if strings.HasPrefix(pinned, cfg.Context.String()) {
		// remove "$(cfg.Context)/" prefix from path to pinned object
		pinned = strings.TrimLeft(pinned, cfg.Context.String()+"/")
	}

	if remove {
		if err := cfg.Pinned.Remove(pinned); err != nil {
			return err
		}
		if err = cfg.DumpToFile(cfgPath); err != nil {
			return err
		}
		output.Minus.Printf("unpinned %v (in '%v').\n", pinned, cfg.Context.String())

	} else {
		if err = cfg.Pinned.Append(pinned); err != nil {
			return err
		}
		if err = cfg.DumpToFile(cfgPath); err != nil {
			return err
		}
		output.Plus.Printf("pinned %v (in '%v').\n", pinned, cfg.Context.String())
	}

	return nil
}
