package app

import (
	"github.com/jieggii/dcfg/internal/commands"
	"github.com/urfave/cli/v2"
	"os"
)

// program metadata
const version = "0.2.0"

// defaults
const defaultConfigFilename = "dcfg.json"

// categories
const filesystemOperationsCategory = "\n  FILESYSTEM OPERATIONS"
const serviceCategory = "\n  SERVICE"
const metadataManagementCategory = "\n  METADATA MANAGEMENT"

func NewApp() *cli.App {
	return &cli.App{
		Name:        "dcfg",
		Usage:       "distribute config",
		UsageText:   "dcfg [--config PATH] <command> [command options]",
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
				Description:  "shows context directory, defined bindings, pinned directories and additions",
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
				Usage:       "append addition",
				UsageText:   "dcfg [--config] add <PATH>",
				Description: "appends new addition",
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
				Category:     metadataManagementCategory,
			},
			{
				Name:        "pin",
				Aliases:     []string{"p"},
				Usage:       "pin or unpin object",
				UsageText:   "dcfg [--config] pin [--remove] <PATH>",
				Description: "pins or unpins file or directory inside context directory",
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
				Category:     metadataManagementCategory,
			},
			// filesystem operations
			{
				Name:         "collect",
				Aliases:      []string{"c"},
				Usage:        "collect additions",
				UsageText:    "dcfg [--config] collect",
				Description:  "copies all registered additions to context directory according to bindings",
				Action:       explicitArgsCountMiddleware(0, commands.Collect),
				OnUsageError: handleUsageError,
				Category:     filesystemOperationsCategory,
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
