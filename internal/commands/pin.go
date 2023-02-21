package commands

import (
	"fmt"
	"github.com/jieggii/dcfg/internal/config"
	"github.com/jieggii/dcfg/internal/output"
	"github.com/urfave/cli/v2"
	"path"
)

func Pin(ctx *cli.Context) error {
	cfgPath := ctx.String("config")
	cfg, err := config.NewConfigFromFile(cfgPath)
	if err != nil {
		return err
	}

	// flags
	remove := ctx.Bool("remove")

	pins := ctx.Args().Slice()

	for _, pin := range pins {
		pin = path.Clean(pin)
		if path.IsAbs(pin) {
			return fmt.Errorf("path to object to be pinned must be relative (got absolute path '%v')", pin)
		}

		if !remove { // if --remove flag is not used
			if err = cfg.Pinned.Append(pin); err != nil {
				output.Error.Println("could not pin '%v' (%v)", pin, err)
				continue
			}
			output.Plus.Printf("pinned '%v'", pin)

		} else { // if --remove flag is used
			if err := cfg.Pinned.Remove(pin); err != nil {
				output.Error.Println("could not unpin '%v' (%v)", pin, err)
				continue
			}
			output.Minus.Printf("unpinned '%v'", pin)
		}
		if err = cfg.DumpToFile(cfgPath); err != nil {
			return err
		}
	}

	return nil
}
