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
	output.Stdout.Printf("dcfg config: '%v'\n", cfgPath)

	// context directory
	output.Stdout.Printf("context directory: '%v'\n", cfg.Context)
	output.Stdout.Println()

	// bindings
	output.Stdout.Println("bindings:")
	if cfg.Bindings.IsPresent() {
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
	if cfg.Additions.IsPresent() {
		for _, addition := range cfg.Additions.Paths {
			// todo: resolve destination path
			// todo: indicate if addition is stored
			destination, resolved := cfg.ResolveAdditionDestination(addition)
			if !resolved {
				destination = "[MISSING SUITABLE BINDING]"
			}

			var collectedString string
			collected, err := cfg.Additions.IsCollected(destination)
			if err != nil {
				collectedString = "???"
				output.Warning.Println(err)
			}
			if collected {
				collectedString = "COLLECTED"
			} else {
				collectedString = "NOT COLLECTED"
			}
			output.Stdout.Printf(" - %-"+strconv.Itoa(cfg.Additions.LongestPathLen)+"v -> %v [%v]", addition, destination, collectedString)
		}
	} else {
		output.Stdout.Println(" * no additions yet *")
	}

	output.Stdout.Println()
	output.Stdout.Println("pinned directories and files:")
	if cfg.Pinned.IsPresent() {
		for _, pinned := range cfg.Pinned {
			output.Stdout.Printf(" - %v\n", pinned)
		}
	} else {
		output.Stdout.Println(" * no pins yet *")
	}
	return nil
}