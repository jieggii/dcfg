package app

import (
	"github.com/jieggii/dcfg/internal/output"
	"github.com/urfave/cli/v2"
	"os"
)

func handleUsageError(ctx *cli.Context, err error, _ bool) error {
	output.Stdout.Printf("usage: %v\n", ctx.Command.UsageText)
	return err
}

func handleCommandNotFoundError(ctx *cli.Context, command string) {
	output.Error.Printf("'%v' is not a dcfg command. See 'dcfg --help'.", command)
	os.Exit(1)
}
