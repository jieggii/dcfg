package commands

import (
	"fmt"
	"github.com/jieggii/dcfg/internal"
	"github.com/jieggii/dcfg/internal/config"
	"github.com/jieggii/dcfg/internal/diff"
	"github.com/jieggii/dcfg/internal/fs"
	"github.com/jieggii/dcfg/internal/input"
	"github.com/jieggii/dcfg/internal/output"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

func Extract(ctx *cli.Context) error {
	cfgPath := ctx.String("config")
	cfg, err := config.NewConfigFromFile(cfgPath)
	if err != nil {
		return err
	}
	nodiff := ctx.Bool("nodiff")

	additionsCount := len(cfg.Additions.Paths)
	if additionsCount == 0 {
		return fmt.Errorf("there are no additions defined")
	}
	output.Stdout.Println("extracting collected additions:")
	for i, addition := range cfg.Additions.Paths {
		destination, found := cfg.ResolveAdditionDestination(addition)
		if !found {
			output.Warning.Printf(
				"cold not resolve destination for '%v'\n",
				addition,
			)
			continue
		}
		var collected bool
		if collected, err = cfg.Additions.IsCollected(destination); err != nil {
			output.Warning.Println(
				"could not check if '%v' collected (%v), skipping\n",
				addition, err,
			)
			continue
		}
		if !collected {
			output.Minus.Printf("skipping uncollected '%v'\n", addition)
			continue
		}
		outputString := fmt.Sprintf(
			"(%v/%v) '%v' -> '%v'",
			i+1, additionsCount, destination, addition,
		)

		output.Stdout.Printf(outputString)

		if !nodiff {
			outputStringLen := len(outputString)
			outputDivider := strings.Repeat("-", outputStringLen)

			output.Stdout.Println(strings.Repeat(" ", int(outputStringLen/2)-4), "diff(s):")
			output.Stdout.Println(outputDivider)

			diffOutput, err := diff.GetDiff(addition, destination)
			if diffOutput != "" {
				output.Stdout.Println(diffOutput)
			}
			output.Stdout.Println(outputDivider)

			if err != nil {
				output.Error.Println("error running diff: %v", err)
			}
		}
		confirm, err := input.ConfirmationPrompt("proceed with copying?")
		if err != nil {
			output.Error.Println(err)
			os.Exit(internal.GenericErrorExitCode)
		}
		if !confirm {
			continue
		}

		if err := fs.Copy(destination, addition); err != nil {
			output.Warning.Println(
				"could not copy '%v' -> '%v' (%v)",
				destination, addition, err,
			)
		}
	}

	return nil
}
