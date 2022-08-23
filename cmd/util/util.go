package util

import (
	"errors"
	"github.com/jieggii/dcfg/cmd/log"
	"github.com/urfave/cli/v2"
	"os"
	"os/user"
	p "path"
	"strings"
)

func expandUser(path string) string {
	if strings.Contains(path, "~") {
		currentUser, err := user.Current()
		if err != nil {
			panic(err) // todo
		}
		return strings.ReplaceAll(path, "~", currentUser.HomeDir)
	}
	return path
}

func CompilePath(path string) string {
	path = p.Clean(path)
	path = expandUser(path)
	return path
}

func LogDryRun() {
	log.Info("(Dry run)")
}

func LogCommandUsageError(context *cli.Context, err error) {
	var commandName string
	var usageText string
	if context.Command.Name != "" {
		usageText = context.Command.UsageText
		commandName = "dcfg " + context.Command.Name
	} else {
		usageText = context.App.UsageText
		commandName = context.App.Name
	}
	log.Error("Error: %v.", err)
	log.Info("Usage: %v.", usageText)
	log.Info("Type `%v --help` for more information.", commandName)
}

func ThrowUsageError(context *cli.Context, message string) {
	if err := HandleUsageError(context, errors.New(message), false); err != nil {
		panic(err)
	}
}

func HandleUsageError(context *cli.Context, err error, _ bool) error {
	LogCommandUsageError(context, err)
	os.Exit(1)
	return nil
}

func CheckCommandArgsCount(context *cli.Context, expected int) {
	if context.Args().Len() != expected {
		ThrowUsageError(context, "too many arguments")
	}
}
