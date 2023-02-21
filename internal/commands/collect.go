package commands

import (
	"github.com/jieggii/dcfg/internal/config"
	"github.com/jieggii/dcfg/internal/fs"
	"github.com/jieggii/dcfg/internal/output"
	"github.com/urfave/cli/v2"
	"strconv"
)

func Collect(ctx *cli.Context) error {
	cfgPath := ctx.String("config")

	cfg, err := config.NewConfigFromFile(cfgPath)
	if err != nil {
		return err
	}

	minAdditionLenStr := strconv.Itoa(cfg.Targets.LongestPathLen)
	for _, target := range cfg.Targets.Paths {
		destination, resolved := cfg.ResolveTargetDestination(target)
		if !resolved {
			output.Warning.Printf(
				"%-"+minAdditionLenStr+"v  : could not resolve target destination (missing suitable binding)",
				target,
			)
			continue
		}
		if err := fs.Copy(target, destination); err != nil {
			output.Warning.Printf(
				"%-"+minAdditionLenStr+"v  : could not copy '%v' to '%v' (%v)",
				target, target, destination, err,
			)
			continue
		}
		output.Plus.Printf("%-"+minAdditionLenStr+"v -> %v", target, destination)
	}
	return nil
}
