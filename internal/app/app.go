package app

import (
	"os"

	"github.com/jieggii/dcfg/internal/commands"
	"github.com/urfave/cli/v2"
)

// dcfg metadata
const version = "0.2.0"

// defaults
const defaultConfigFilename = "dcfg.json"
const defaultDiffBinPath = "/usr/bin/diff"

// categories
const serviceCategory = "SERVICE"
const metadataManagementCategory = "METADATA MANAGEMENT"
const filesystemOperationsCategory = "FILESYSTEM OPERATIONS"

func NewApp() *cli.App {
	return &cli.App{
		Name:        "dcfg",
		Usage:       "distribute config",
		UsageText:   "dcfg [--help] [--version] [--config PATH] <command> [command options]",
		Version:     version,
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
			// service
			{
				Name:         "init",
				Aliases:      []string{"i"},
				Usage:        "initialize dcfg",
				UsageText:    "dcfg [--config PATH] init",
				Description:  "creates dcfg config file",
				Action:       explicitArgsCountMiddleware(0, commands.Init),
				OnUsageError: handleUsageError,
				Category:     serviceCategory,
			},
			{
				Name:         "status",
				Aliases:      []string{"s"},
				Usage:        "show information about current state",
				UsageText:    "dcfg [--config] status",
				Description:  "shows useful information about targets, defined bindings and pinned nodes",
				Action:       explicitArgsCountMiddleware(0, commands.Status),
				OnUsageError: handleUsageError,
				Category:     serviceCategory,
			},
			// metadata management
			{
				Name:        "bind",
				Aliases:     []string{"b"},
				Usage:       "register or remove binding",
				UsageText:   "dcfg [--config PATH] bind [--remove] <source> [destination]",
				Description: "registers or removes absolute-to-relative path binding",
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
				Category:     metadataManagementCategory,
			},
			{
				Name:        "add",
				Aliases:     []string{"a"},
				Usage:       "add target(s)",
				UsageText:   "dcfg [--config] add <TARGET1> ...",
				Description: "adds new target(s)",
				Action:      moreThanArgsCountMiddleware(1, commands.Add),
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "collect",
						Usage:   "collect added target(s)",
						Value:   false,
						Aliases: []string{"c"},
					},
				},
				OnUsageError: handleUsageError,
				Category:     metadataManagementCategory,
			},
			{
				Name:        "pin",
				Aliases:     []string{"p"},
				Usage:       "pin or unpin object",
				UsageText:   "dcfg [--config] pin [--remove] <PATH1> ...",
				Description: "pins or unpins file or directory",
				Action:      moreThanArgsCountMiddleware(1, commands.Pin),
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "remove",
						Value:   false,
						Usage:   "remove pinned object",
						Aliases: []string{"r"},
					},
				},
				OnUsageError: handleUsageError,
				Category:     metadataManagementCategory,
			},
			// filesystem operations
			{
				Name:         "collect",
				Aliases:      []string{"c"},
				Usage:        "collect targets",
				UsageText:    "dcfg [--config] collect",
				Description:  "copies all targets according to the bindings",
				Action:       explicitArgsCountMiddleware(0, commands.Collect),
				OnUsageError: handleUsageError,
				Category:     filesystemOperationsCategory,
			},
			{
				Name:        "extract",
				Aliases:     []string{"e"},
				Usage:       "extract collected targets",
				UsageText:   "dcfg [--config] extract [--hard] [--no-diff] [--overwrite-source-prefix oldPrefix:newPrefix]",
				Description: "copies all collected targets to their sources",
				Action:      explicitArgsCountMiddleware(0, commands.Extract),
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "no-diff",
						Value: false,
						Usage: "do not show diff(s)",
					},
					&cli.BoolFlag{
						Name:  "hard",
						Value: false,
						Usage: "disable confirmation prompts before performing operations (dangerous)",
					},
					&cli.StringSliceFlag{
						Name:    "overwrite-source-prefix",
						Aliases: []string{"o"},
						Usage:   "overwrites target source prefix this time",
					},
					// hidden flags:
					&cli.StringFlag{
						Name:   "diff-bin-path",
						Hidden: true,
						Value:  defaultDiffBinPath,
					},
				},
				OnUsageError: handleUsageError,
				Category:     filesystemOperationsCategory,
			},
			{
				Name:        "remove",
				Aliases:     []string{"rm"},
				Usage:       "remove target",
				UsageText:   "dcfg [--config] remove <TARGET1> ...",
				Description: "removes target",
				Action:      moreThanArgsCountMiddleware(1, commands.Remove),
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "soft",
						Usage:   "do not remove stored (collected) targets",
						Value:   false,
						Aliases: []string{"s"},
					},
				},
				OnUsageError: handleUsageError,
				Category:     filesystemOperationsCategory,
			},
			{
				Name:        "clean",
				Aliases:     []string{"cl"},
				Usage:       "remove all outdated collected targets and other trash",
				UsageText:   "dcfg [--config] clean",
				Description: "removes everything except up-to-date collected targets, pinned nodes and some useful files",
				Action:      commands.Clean,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "yes",
						Usage:   "disable confirmation prompt (dangerous)",
						Value:   false,
						Aliases: []string{"y"},
					},
				},
				OnUsageError: handleUsageError,
				Category:     filesystemOperationsCategory,
			},
		},
		HideHelpCommand: true,
		CommandNotFound: handleCommandNotFoundError,
		OnUsageError:    handleUsageError,
		Authors: []*cli.Author{
			{Name: "jieggii", Email: "jieggii@protonmail.com"},
		},

		BashComplete: cli.DefaultAppComplete,
		Reader:       os.Stdin,
		Writer:       os.Stdout,
		ErrWriter:    os.Stderr,
	}
}
