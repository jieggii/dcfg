package main

import (
	"fmt"
	"github.com/jieggii/dcfg/cmd/commands"
	"github.com/jieggii/dcfg/cmd/util"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := &cli.App{
		Name:        "dcfg",
		Usage:       "distribute config",
		UsageText:   "dcfg [global options] [command] [command options]",
		Version:     "0.1.0",
		Description: "Simple tool for copying and distributing your config files.",
		Authors: []*cli.Author{
			{
				Name:  "jieggii",
				Email: "jieggii@pm.me",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Usage:   "dcfg config file `PATH`",
				Value:   "dcfg.conf",
				Aliases: []string{"c"},
			},
			&cli.BoolFlag{
				Name:    "dry",
				Usage:   "dry run",
				Value:   false,
				Aliases: []string{"d"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:         "init",
				Aliases:      []string{"i"},
				Usage:        "create dcfg config file",
				UsageText:    "dcfg init",
				Action:       commands.Init,
				OnUsageError: util.HandleUsageError,
			},
			{
				Name:         "add",
				Aliases:      []string{"a"},
				Usage:        "copy additions to the context directory",
				UsageText:    "dcfg add",
				Action:       commands.Add,
				OnUsageError: util.HandleUsageError,
			},
			{
				Name:         "clean",
				Aliases:      []string{"c"},
				Usage:        "remove everything from context directory excepting additions, pins, dcfg config file and .git directory",
				UsageText:    "dcfg clean",
				Action:       commands.Clean,
				OnUsageError: util.HandleUsageError,
			},
			{
				Name:         "push",
				Aliases:      []string{"p"},
				Usage:        "add, commit and push changes to a remote git repository",
				UsageText:    "dcfg push",
				Action:       commands.Push,
				OnUsageError: util.HandleUsageError,
			},
		},
		HideHelpCommand: true,
		Action: func(context *cli.Context) error {
			args := context.Args()
			if args.Len() == 0 {
				util.ThrowUsageError(context, "no command provided")
			} else {
				util.ThrowUsageError(context, fmt.Sprintf("unknown command `%v`", args.First()))
			}
			os.Exit(1)
			return nil
		},
		OnUsageError: util.HandleUsageError,
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Error: %v.", err)
	}
}
