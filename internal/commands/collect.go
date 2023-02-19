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

	minAdditionLenStr := strconv.Itoa(cfg.Additions.LongestPathLen)
	for _, addition := range cfg.Additions.Paths {
		destination, resolved := cfg.ResolveAdditionDestination(addition)
		if !resolved {
			output.Warning.Printf(
				"%-"+minAdditionLenStr+"v  : could not resolve addition destination (missing suitable binding)",
				addition,
			)
			continue
		}
		if err := fs.Copy(addition, destination); err != nil {
			output.Warning.Printf(
				"%-"+minAdditionLenStr+"v  : could not copy '%v' to '%v' (%v)",
				addition, addition, destination, err,
			)
			continue
		}
		output.Plus.Printf("%-"+minAdditionLenStr+"v -> %v", addition, destination)
	}
	return nil
}
