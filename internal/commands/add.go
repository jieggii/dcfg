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

const missingSuitableBindingErrorText = "could not resolve addition destination (missing suitable binding)"

func Add(ctx *cli.Context) error {
	cfgPath := ctx.String("config")
	cfg, err := config.NewConfigFromFile(cfgPath)
	if err != nil {
		return err
	}

	//verbose := ctx.Bool("verbose")
	collect := ctx.Bool("collect")

	additions := ctx.Args().Slice()

	destinations := map[string]string{}

	// at first validating:
	for _, addition := range additions {
		addition = path.Clean(addition)
		if !path.IsAbs(addition) {
			return fmt.Errorf("path to addition must be absolute (got relative path '%v')", addition)
		}
		destination, resolved := cfg.ResolveAdditionDestination(addition)
		if !resolved {
			return errors.New(missingSuitableBindingErrorText)
		}
		destinations[addition] = destination
	}

	// and only then performing actions:
	for _, addition := range additions {
		addition = path.Clean(addition)
		destination := destinations[addition]

		if err := cfg.Additions.Append(addition); err != nil {
			output.Error.Println(
				"could not append '%v' to additions (%v)", addition, err,
			)
			continue
		}
		output.Plus.Printf("appended new addition '%v'", addition)

		//if verbose {
		//	output.Verbose.Printf("the addition will be copied to '%v' when collecting", destination)
		//}

		if collect {
			if err := fs.Copy(addition, destination); err != nil {
				output.Error.Printf(
					"could not copy '%v' to '%v' (%v)", addition, destination, err,
				)
				continue
			}
			output.Plus.Printf("copied '%v' -> '%v'", addition, destination)
		}
	}

	if err := cfg.DumpToFile(cfgPath); err != nil {
		return err
	}

	return nil
}
