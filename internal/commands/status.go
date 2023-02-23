package commands

import (
	"github.com/jieggii/dcfg/internal/config"
	"github.com/jieggii/dcfg/internal/output"
	"github.com/urfave/cli/v2"
	"strconv"
)

// Status command outputs some useful information about current state.
func Status(ctx *cli.Context) error {
	cfgPath := ctx.String("config")
	cfg, err := config.NewConfigFromFile(cfgPath)
	if err != nil {
		return err
	}

	// config file
	output.Stdout.Printf("dcfg config file: '%v'", cfgPath)

	// bindings
	output.Stdout.Println("bindings:")
	if cfg.Bindings.Any() {
		for i, source := range cfg.Bindings.Sources {
			destination, err := cfg.Bindings.ResolveDestination(source)
			if err != nil {
				panic(err)
			}
			output.Stdout.Printf(" %v. %-"+strconv.Itoa(cfg.Bindings.LongestSourceLen)+"v -> %v", i+1, source, destination)
		}
	} else {
		output.Stdout.Println(" * no bindings yet *")
	}
	output.Stdout.Println()

	// targets
	output.Stdout.Println("targets:")
	if cfg.Targets.Any() {
		longestTargetPathLenString := strconv.Itoa(cfg.Targets.LongestPathLen)
		for _, target := range cfg.Targets.Paths {
			destination, resolved := cfg.ResolveTargetDestination(target)

			var collectedIndicator string
			if resolved { // if destination was resolved
				collected, err := cfg.Targets.IsCollected(destination)
				if collected { // target is collected
					collectedIndicator = "+"
				} else if !collected && err == nil { // target is not collected
					collectedIndicator = "-"
				} else { // could not check if target is collected
					collectedIndicator = "?"
					output.Warning.Println(err)
				}
			} else { // destination was not resolved (missing suitable binding)
				collectedIndicator = "!"
				destination = "[MISSING SUITABLE BINDING]"
			}
			output.Stdout.Printf(" %v %-"+longestTargetPathLenString+"v -> %v", collectedIndicator, target, destination)
		}
	} else {
		output.Stdout.Println(" * no targets yet *")
	}

	output.Stdout.Println()
	output.Stdout.Println("pinned objects:")
	if cfg.Pinned.Any() {
		for _, pinned := range cfg.Pinned {
			output.Stdout.Printf(" - %v", pinned)
		}
	} else {
		output.Stdout.Println(" * no pins yet *")
	}
	return nil
}
