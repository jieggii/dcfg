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

	sourcePrefixReplaceMap := make(map[string]string)

	// fill sourcePrefixReplaceMap - map of prefix replacements of targets
	// old -> new
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

		// check if old prefix exists in any of targets:
		oldPrefixIsPresent := false
		for _, target := range cfg.Targets.Paths {
			if strings.HasPrefix(target, oldPrefix) {
				oldPrefixIsPresent = true
				break
			}
		}

		if !oldPrefixIsPresent {
			return fmt.Errorf(
				"there is not prefix '%v' in any of the targets",
				oldPrefix,
			)
		}

		sourcePrefixReplaceMap[oldPrefix] = newPrefix
	}

	// print sourcePrefixReplaceMap
	//output.Stdout.Println(
	//	"overwritten target destinations:",
	//)
	//for oldPrefix, newPrefix := range sourcePrefixReplaceMap {
	//	output.Stdout.Printf("%v -> %)
	//}

	targetsCount := len(cfg.Targets.Paths)
	if targetsCount == 0 {
		return fmt.Errorf("there are no targets")
	}
	for i, target := range cfg.Targets.Paths {
		if i > 0 && i < targetsCount-1 {
			output.Stdout.Println() // add blank line for more pretty output
		}
		destination, found := cfg.ResolveTargetDestination(target)
		if !found {
			output.Warning.Printf(
				"cold not resolve destination for '%v'",
				target,
			)
			continue
		}
		// overwrite target prefix
		for oldPrefix, newPrefix := range sourcePrefixReplaceMap {
			if strings.HasPrefix(target, oldPrefix) {
				newTarget := strings.Replace(target, oldPrefix, newPrefix, 1)
				output.Warning.Printf(
					"using '%v' instead of '%v'",
					newTarget, target,
				)
				target = newTarget
			}
		}

		var collected bool
		if collected, err = cfg.Targets.IsCollected(destination); err != nil {
			output.Warning.Println(
				"could not check if '%v' collected (%v), skipping",
				target, err,
			)
			continue
		}
		if !collected { // if target is not collected
			output.Minus.Printf("skipping uncollected '%v'", target)
			continue
		}

		outputString := fmt.Sprintf(
			"Target %v/%v: '%v' -> '%v'",
			i+1, targetsCount, destination, target,
		)
		output.Stdout.Printf(outputString)

		if !nodiff { // if --no-diff flag is not used
			outputStringLen := len(outputString)
			outputDivider := strings.Repeat("-", outputStringLen)

			output.Stdout.Println(strings.Repeat(" ", int(outputStringLen/2)-4), "diff(s):")
			output.Stdout.Println(outputDivider)

			diffOutput, err := diff.GetDiff(diffBinPath, target, destination)
			if diffOutput != "" {
				output.Stdout.Println(diffOutput)
			}
			output.Stdout.Println(outputDivider)

			if err != nil {
				output.Error.Printf("error running diff (%v)", err)
			}
		}
		if !hard { // if --hard flag is not used
			confirm, err := input.ConfirmationPrompt("proceed with copying?")
			if err != nil {
				return err
			}
			if !confirm {
				output.Warning.Printf("skipping '%v'", target)
				continue
			}
		}

		if err := fs.Copy(destination, target); err != nil {
			output.Warning.Println(
				"could not copy '%v' -> '%v' (%v)",
				destination, target, err,
			)
			output.Warning.Printf("skipping '%v'", target)
			continue
		}
		output.Plus.Printf("'%v' -> '%v'", destination, target)
	}

	return nil
}
