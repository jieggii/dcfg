package app

import (
	"github.com/jieggii/dcfg/internal"
	"github.com/jieggii/dcfg/internal/output"
	"github.com/urfave/cli/v2"
	"os"
)

func handleUsageError(ctx *cli.Context, err error, _ bool) error {
	output.Error.Println(err)
	output.Stdout.Printf("usage: %v", ctx.Command.UsageText)
	os.Exit(internal.GenericErrorExitCode)
	return nil
}

func handleCommandNotFoundError(_ *cli.Context, command string) {
	output.Error.Printf("'%v' is not a dcfg command. See 'dcfg --help'", command)
	os.Exit(internal.GenericErrorExitCode)
}
