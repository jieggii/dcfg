package commands

import (
	"fmt"
	"github.com/jieggii/dcfg/internal/config"
	"github.com/jieggii/dcfg/internal/diff"
	"github.com/jieggii/dcfg/internal/fs"
	"github.com/jieggii/dcfg/internal/input"
	"github.com/jieggii/dcfg/internal/output"
	"github.com/urfave/cli/v2"
	"path"
	"strings"
)

func Extract(ctx *cli.Context) error {
	cfgPath := ctx.String("config")
	cfg, err := config.NewConfigFromFile(cfgPath)
	if err != nil {
		return err
	}

	// options
	nodiff := ctx.Bool("nodiff")
	hard := ctx.Bool("hard")
	sourcePrefixReplacements := ctx.StringSlice("overwrite-source-prefix")

	// hidden options
	diffBinPath := ctx.String("diff-bin-path")

	sourcePrefixReplaceMap := map[string]string{}

	// fill sourcePrefixReplaceMap
	for _, directive := range sourcePrefixReplacements {
		tokens := strings.Split(directive, ":")
		if len(tokens) != 2 {
			return fmt.Errorf(
				"invalid directive '%v' for --overwrite-source-prefix option",
				directive,
			)
		}

		oldPrefix := path.Clean(tokens[0])
		newPrefix := path.Clean(tokens[1])

		// check if prefix (being replaced) exists in any of additions:
		oldPrefixIsPresent := false
		for _, addition := range cfg.Additions.Paths {
			if strings.HasPrefix(addition, oldPrefix) {
				oldPrefixIsPresent = true
				break
			}
		}

		if !oldPrefixIsPresent {
			return fmt.Errorf("there is not prefix '%v' in any of the additions", oldPrefix)
		}

		sourcePrefixReplaceMap[oldPrefix] = newPrefix
	}

	fmt.Println(sourcePrefixReplaceMap)
	additionsCount := len(cfg.Additions.Paths)
	if additionsCount == 0 {
		return fmt.Errorf("there are no additions defined")
	}
	for i, addition := range cfg.Additions.Paths {
		if i > 0 && i < additionsCount-1 {
			output.Stdout.Println() // add blank line for more pretty output
		}
		destination, found := cfg.ResolveAdditionDestination(addition)
		if !found {
			output.Warning.Printf(
				"cold not resolve destination for '%v'",
				addition,
			)
			continue
		}
		var collected bool
		if collected, err = cfg.Additions.IsCollected(destination); err != nil {
			output.Warning.Println(
				"could not check if '%v' collected (%v), skipping",
				addition, err,
			)
			continue
		}
		if !collected { // if addition is not collected
			output.Minus.Printf("skipping uncollected '%v'", addition)
			continue
		}

		outputString := fmt.Sprintf(
			"Addition %v/%v: '%v' -> '%v'",
			i+1, additionsCount, destination, addition,
		)
		output.Stdout.Printf(outputString)

		if !nodiff { // if --no-diff flag is not used
			outputStringLen := len(outputString)
			outputDivider := strings.Repeat("-", outputStringLen)

			output.Stdout.Println(strings.Repeat(" ", int(outputStringLen/2)-4), "diff(s):")
			output.Stdout.Println(outputDivider)

			diffOutput, err := diff.GetDiff(diffBinPath, addition, destination)
			if diffOutput != "" {
				output.Stdout.Println(diffOutput)
			}
			output.Stdout.Println(outputDivider)

			if err != nil {
				output.Error.Println("error running diff: %v", err)
			}
		}
		if !hard { // if --hard flag is not used
			confirm, err := input.ConfirmationPrompt("proceed with copying?")
			if err != nil {
				return err
			}
			if !confirm {
				output.Warning.Printf("skipping '%v'", addition)
				continue
			}
		}

		if err := fs.Copy(destination, addition); err != nil {
			output.Warning.Println(
				"could not copy '%v' -> '%v' (%v)",
				destination, addition, err,
			)
			output.Warning.Printf("skipping '%v'", addition)
			continue
		}
		output.Plus.Printf("'%v' -> '%v'", destination, addition)
	}

	return nil
}
