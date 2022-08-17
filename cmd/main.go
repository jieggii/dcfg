package main

import (
	"fmt"
	"github.com/jieggii/dcfg/cmd/config"
	cp "github.com/otiai10/copy"
	"github.com/urfave/cli/v2"
	"os"
	"strconv"
	"strings"
)

func usageErrorHandler(context *cli.Context, err error, isSubCommand bool) error {
	// todo: add command name (init/update if context command is init/update)
	fmt.Printf("Error: %v.\n", err)
	fmt.Println("Type `dcfg --help` for more information.")
	os.Exit(1)
	return nil
}

func logCommandUsageError(context *cli.Context, err string) {
	var commandName string
	var usageText string
	if context.Command.Name != "" {
		usageText = context.Command.UsageText
		commandName = "dcfg " + context.Command.Name
	} else {
		usageText = context.App.UsageText
		commandName = context.App.Name
	}
	fmt.Printf("Error: %v.\n", err)
	fmt.Printf("Usage: %v.\n", usageText)
	fmt.Printf("Type `%v --help` for more information.\n", commandName)
}

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
		},
		Commands: []*cli.Command{
			{
				Name:      "init",
				Aliases:   []string{"i"},
				Usage:     "create dcfg config file in the current working directory",
				UsageText: "dcfg init",
				Action: func(context *cli.Context) error {
					cfgPath := context.String("config")
					if context.Args().Len() != 0 {
						logCommandUsageError(context, "too many arguments")
						os.Exit(1)
					}
					if err := config.CreateConfig(cfgPath); err != nil {
						fmt.Printf("Error: could not create dcfg config file `%v`: %v.\n", cfgPath, err)
						os.Exit(2)
					}
					fmt.Printf("Created dcfg config file `%v` in the current working directory.\n", cfgPath)
					return nil
				},
				OnUsageError: usageErrorHandler,
			},
			{
				Name:      "update",
				Aliases:   []string{"u"},
				Usage:     "update current working directory according to the dcfg config file",
				UsageText: "dcfg update",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "remove",
						Aliases: []string{"r"},
						Usage:   "remove directories and files that are not declared in config file as additions",
						Value:   false,
					},
				},
				Action: func(context *cli.Context) error {
					cfgPath := context.String("config")
					cfg := config.ReadConfig(cfgPath)
					if len(cfg.Additions.Paths) == 0 {
						fmt.Printf("Error: nothing to do (there are no `add` directives in `%v`).\n", cfgPath)
						return nil
					}
					fmt.Println("Updating working directory according to these bindings:")
					for i, source := range cfg.Bindings.Sources {
						destination := cfg.Bindings.Destinations[i]
						fmt.Printf("%v. %-"+strconv.Itoa(cfg.Bindings.LongestSourceLength)+"v : %v\n", i+1, source, destination)
					}
					fmt.Println("")
					for _, globalPath := range cfg.Additions.Paths {
						matched := false
						for i, source := range cfg.Bindings.Sources {
							destination := cfg.Bindings.Destinations[i]
							if strings.HasPrefix(globalPath, source) {
								matched = true
								localPath := strings.Replace(globalPath, source, destination, 1)
								if err := cp.Copy(globalPath, localPath); err != nil {
									fmt.Printf("Error: could not copy %v to %v\n", globalPath, localPath)
								} else {
									fmt.Printf("%-"+strconv.Itoa(cfg.Additions.LongestPathLength)+"v -> %v\n", globalPath, localPath)
								}
								break
							}
						}
						if !matched {
							fmt.Printf(
								"Warning: ingoring unmatched addition: %v\n",
								globalPath,
							)
						}
					}

					return nil
				},
				OnUsageError: usageErrorHandler,
			},
			{
				Name:      "push",
				Aliases:   []string{"p"},
				Usage:     "add, commit and push changes to a remote git repository",
				UsageText: "dcfg push",
				Action: func(context *cli.Context) error {
					return nil
				},
				OnUsageError: usageErrorHandler,
			},
		},
		HideHelpCommand: true,
		Action: func(context *cli.Context) error {
			args := context.Args()
			if args.Len() == 0 {
				logCommandUsageError(context, "no command provided")
			} else {
				logCommandUsageError(context, fmt.Sprintf("unknown command `%v`", args.First()))
			}
			os.Exit(1)
			return nil
		},
		OnUsageError: usageErrorHandler,
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Error: %v.", err)
	}
}
