package main

import (
	"fmt"
	"github.com/jieggii/dcfg/internal/commands"
	"github.com/jieggii/dcfg/internal/output"
	"github.com/urfave/cli/v2"
	"os"
)

const defaultConfigFilename = "dcfg.json"

func intervalArgsCountMiddleware(minArgsCount int, maxArgsCount int, action cli.ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		args := ctx.Args()
		argsCount := args.Len()
		if argsCount >= minArgsCount && argsCount <= maxArgsCount {
			return action(ctx)
		} else {
			return fmt.Errorf(
				"%v command takes from %v to %v arguments, got %v.\nusage: %v",
				ctx.Command.Name, minArgsCount, maxArgsCount, argsCount, ctx.Command.UsageText,
			)
		}
	}
}

func explicitArgsCountMiddleware(expectedArgsCount int, action cli.ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		args := ctx.Args()
		if args.Len() == expectedArgsCount {
			return action(ctx)
		} else {
			if expectedArgsCount == 0 {
				return fmt.Errorf(
					"%v command takes no arguments, got %v.\nusage: %v",
					ctx.Command.Name, args.Len(), ctx.Command.UsageText,
				)
			} else {
				return fmt.Errorf(
					"%v command takes exactly %v argument(s), got %v.\nusage: %v",
					ctx.Command.Name, expectedArgsCount, args.Len(), ctx.Command.UsageText,
				)
			}
		}
	}
}

func handleUsageError(ctx *cli.Context, err error, _ bool) error {
	output.Stdout.Printf("usage: %v\n", ctx.Command.UsageText)
	return err
}

func main() {
	app := &cli.App{
		Name:        "dcfg",
		Usage:       "distribute config",
		UsageText:   "dcfg [--config PATH] command [command options]",
		Version:     "0.2.0",
		Description: "Minimalist tool for copying, storing and distributing your system-wide and user config files.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Usage:   "dcfg config file `PATH`",
				Value:   defaultConfigFilename,
				Aliases: []string{"c"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:         "init",
				Aliases:      []string{"i"},
				Usage:        "initialize dcfg",
				UsageText:    "dcfg [--config PATH] init",
				Description:  "creates dcfg config file",
				Action:       explicitArgsCountMiddleware(0, commands.Init),
				OnUsageError: handleUsageError,
			},
			{
				Name:        "bind",
				Aliases:     []string{"b"},
				Usage:       "create new binding",
				UsageText:   "dcfg [--config PATH] bind [--remove] <source> [destination]",
				Description: "creates new path-to-path binding (absolute to relative)",
				Action:      intervalArgsCountMiddleware(1, 2, commands.Bind),
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "remove",
						Usage:   "remove binding",
						Value:   false,
						Aliases: []string{"r"},
					},
				},
				OnUsageError: handleUsageError,
			},

			{
				Name:        "add",
				Aliases:     []string{"a"},
				Usage:       "append addition",
				UsageText:   "dcfg [--config] add <PATH>",
				Description: "creates new addition (without copying it)",
				Action:      explicitArgsCountMiddleware(1, commands.Add),
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "collect",
						Usage:   "copy addition to context directory according to the bindings",
						Value:   false,
						Aliases: []string{"c"},
					},
				},
				OnUsageError: handleUsageError,
			},
			{
				Name:        "remove",
				Aliases:     []string{"rm"},
				Usage:       "remove addition",
				UsageText:   "dcfg [--config] remove <PATH>",
				Description: "removes addition",
				Action:      explicitArgsCountMiddleware(1, commands.Remove),
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "soft",
						Usage:   "do not remove collected objects from context directory",
						Value:   false,
						Aliases: []string{"s"},
					},
				},
				OnUsageError: handleUsageError,
			},
			{
				Name:        "pin",
				Aliases:     []string{"p"},
				Usage:       "pin file or directory",
				UsageText:   "dcfg [--config] pin [--remove] <PATH>",
				Description: "pins file or directory inside context directory",
				Action:      explicitArgsCountMiddleware(1, commands.Pin),
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "remove",
						Value:   false,
						Usage:   "remove pinned object",
						Aliases: []string{"r"},
					},
				},
				OnUsageError: handleUsageError,
			},
			{
				Name:         "status",
				Aliases:      []string{"s"},
				Usage:        "shows information about current state",
				UsageText:    "dcfg [--config] status",
				Description:  "shows context directory, defined bindings, pinned directories and additions",
				Action:       explicitArgsCountMiddleware(0, commands.Status),
				OnUsageError: handleUsageError,
			},
		},
		HideHelpCommand: true,
		CommandNotFound: func(ctx *cli.Context, command string) {
			output.Error.Printf("'%v' is not a dcfg command. See 'dcfg --help'.", command)
			os.Exit(1)
		},
		OnUsageError: handleUsageError,
		Authors: []*cli.Author{
			{Name: "jieggii", Email: "jieggii@pm.me"},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		output.Error.Printf("%v.\n", err)
		os.Exit(1)
	}
}
