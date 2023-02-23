package commands

import (
	"github.com/jieggii/dcfg/internal/config"
	"github.com/jieggii/dcfg/internal/input"
	"github.com/jieggii/dcfg/internal/output"
	"github.com/jieggii/dcfg/internal/util"
	"github.com/urfave/cli/v2"
	"os"
)

// Files and directories which will be not deleted anyway
var ignoredByDefault = []string{".git", ".gitignore"}

// Clean command removes trash.
func Clean(ctx *cli.Context) error {
	cfgPath := ctx.String("config")
	cfg, err := config.NewConfigFromFile(cfgPath)
	if err != nil {
		return err
	}

	// options
	confirmation := ctx.Bool("yes")

	nodes, err := os.ReadDir(".")
	if err != nil {
		return err
	}

	ignoredByDefault = append(ignoredByDefault, cfgPath) // add current dcfg config to ignored by default

	var deletions []string
	for _, node := range nodes {
		name := node.Name()

		if util.ItemIsInArray(ignoredByDefault, name) { // if node is ignored by default
			continue
		}
		if cfg.Pinned.Exists(name) { // if node is pinned
			continue
		}
		if cfg.Bindings.DestinationWithPrefixExists(name) { // if node is binding destination
			continue
		}
		deletions = append(deletions, name)
	}
	if len(deletions) == 0 {
		output.Warning.Println("nothing to delete")
		return nil
	}

	output.Stdout.Printf("nodes to be deleted:")
	for _, path := range deletions {
		output.Stdout.Printf("- '%v'", path)
	}
	if !confirmation {
		output.Stdout.Println()
		confirmation, err = input.ConfirmationPrompt("proceed with deletions?")
		if err != nil {
			return err
		}
	}
	if confirmation {
		for _, name := range deletions {
			if err := os.RemoveAll(name); err != nil {
				output.Warning.Printf("could not delete '%v' (%v)", name, err)
				continue
			}
			output.Minus.Printf("removed '%v'", name)
		}
	}
	return nil
}
