package commands

import (
	"github.com/jieggii/dcfg/internal/config"
	"github.com/jieggii/dcfg/internal/output"
	"github.com/urfave/cli/v2"
	"strconv"
)

func Status(ctx *cli.Context) error {
	cfgPath := ctx.String("config")

	cfg, err := config.NewConfigFromFile(cfgPath)
	if err != nil {
		return err
	}
	// config file
	output.Stdout.Printf("dcfg config file: '%v'\n", cfgPath)

	// bindings
	output.Stdout.Println("bindings:")
	if cfg.Bindings.Any() {
		for i, source := range cfg.Bindings.Sources {
			destination, err := cfg.Bindings.ResolveDestination(source)
			if err != nil {
				panic(err)
			}
			output.Stdout.Printf(" %v. %-"+strconv.Itoa(cfg.Bindings.LongestSourceLen)+"v -> %v\n", i+1, source, destination)
		}
	} else {
		output.Stdout.Println(" * no bindings yet *")
	}
	output.Stdout.Println()

	// additions
	output.Stdout.Println("additions:")
	if cfg.Additions.Any() {
		longestAdditionLenString := strconv.Itoa(cfg.Additions.LongestPathLen)
		for _, addition := range cfg.Additions.Paths {
			destination, resolved := cfg.ResolveAdditionDestination(addition)

			var collectedIndicator string
			if resolved { // if destination was resolved
				collected, err := cfg.Additions.IsCollected(destination)
				if collected { // addition is collected
					collectedIndicator = "+"
				} else if !collected && err == nil { // addition is not collected
					collectedIndicator = "-"
				} else { // could not check if addition is collected
					collectedIndicator = "?"
					output.Warning.Println(err)
				}
			} else { // destination was not resolved (missing suitable binding)
				collectedIndicator = "!"
				destination = "[MISSING SUITABLE BINDING]"
			}
			output.Stdout.Printf(" %v %-"+longestAdditionLenString+"v -> %v", collectedIndicator, addition, destination)
		}
	} else {
		output.Stdout.Println(" * no additions yet *")
	}

	output.Stdout.Println()
	output.Stdout.Println("pinned objects:")
	if cfg.Pinned.Any() {
		for _, pinned := range cfg.Pinned {
			output.Stdout.Printf(" - %v\n", pinned)
		}
	} else {
		output.Stdout.Println(" * no pins yet *")
	}
	return nil
}
