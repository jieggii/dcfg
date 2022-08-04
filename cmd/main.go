package main

import (
	"errors"
	"fmt"
	"github.com/jieggii/dcfg/cmd/config"
	"github.com/urfave/cli/v2"
	"os"
)

func logCommandUsageError(err error, command *cli.Command) {
	fmt.Printf("Error: %v.\n", err)
	fmt.Printf("Usage: %v. Type `dcfg %v --help` for more information.\n", command.UsageText, command.Name)
}

func main() {
	app := &cli.App{
		Name:        "dcfg",
		Usage:       "distribute config",
		UsageText:   "dcfg [command] [flags]",
		Version:     "0.1.0",
		Description: "bruh",
		Authors: []*cli.Author{
			{
				Name:  "jieggii",
				Email: "jieggii@pm.me",
			},
		},
		Commands: []*cli.Command{
			{
				Name:        "init",
				Aliases:     []string{"i"},
				Usage:       "create dcfg config file in the current working dir",
				UsageText:   "dcfg init",
				Description: "Creates dcfg.conf file in the current directory.",
				Action: func(context *cli.Context) error {
					cfgFilename := "dcfg.conf"
					if context.Args().Len() != 0 {
						logCommandUsageError(errors.New("too many arguments"), context.Command)
						os.Exit(1)
					}
					if err := config.CreateConfig(cfgFilename); err != nil {
						fmt.Printf("Error: could not create dcfg config file `%v`: %v.\n", cfgFilename, err)
						os.Exit(1)
					}
					fmt.Printf("Created dcfg config file `%v` in the current working directory.\n", cfgFilename)
					return nil
				},
			},
			{
				Name:      "update",
				Aliases:   []string{"u"},
				Usage:     "update current working dir according to the dcfg config file",
				UsageText: "dcfg update",
				Description: "Updates current working directory according to the dcfg config file. " +
					"Asks about deletions and overwritings if --yes flag is not used.",
				Flags: []cli.Flag{},
				Action: func(context *cli.Context) error {
					cfgFilename := "dcfg.conf"
					cfg, err := config.ReadConfig(cfgFilename)
					if err != nil {
						fmt.Printf("Could not read dcfg config file `%v`: %v.\n", cfgFilename, err)
					}
					fmt.Println(cfg)
					if len(cfg.Additions) == 0 {
						fmt.Printf("Nothing to do (there are no `add` directives in `%v`).\n", cfgFilename)
						return nil
					}
					for _, addition := range cfg.Additions {
						fmt.Println(addition)
					}
					return nil
				},
			},
		},
		HideHelpCommand: true,
		Action: func(context *cli.Context) error {

			return nil
		},
		OnUsageError: func(context *cli.Context, err error, isSubCommand bool) error {
			fmt.Printf("Error: %v.\n", err)
			fmt.Println("Type `dcfg --help` for more information.")
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Error: %v.", err)
	}
}
