package app

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

// intervalArgsCountMiddleware ensures that count of arguments passed to the command
// is in the provided interval [minArgsCount, maxArgsCount].
func intervalArgsCountMiddleware(minArgsCount int, maxArgsCount int, action cli.ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		args := ctx.Args()
		argsCount := args.Len()
		if argsCount >= minArgsCount && argsCount <= maxArgsCount {
			return action(ctx)
		} else {
			return fmt.Errorf(
				"%v command takes from %v to %v arguments, got %v\nusage: %v",
				ctx.Command.Name, minArgsCount, maxArgsCount, argsCount, ctx.Command.UsageText,
			)
		}
	}
}

// atLeastOneArgumentMiddleware ensures that at least one argument was passed to command.
func atLeastOneArgumentMiddleware(action cli.ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		argsCount := ctx.Args().Len()
		if argsCount < 1 {
			return fmt.Errorf(
				"%v command takes at least 1 argument, got %v\nusage: %v",
				ctx.Command.Name, argsCount, ctx.Command.UsageText,
			)
		}
		return action(ctx)
	}
}

// noArgsMiddleware ensures that no arguments were passed to command.
func noArgsMiddleware(action cli.ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		argsCount := ctx.Args().Len()
		if argsCount != 0 {
			return fmt.Errorf(
				"%v command takes no arguments, got %v\nusage: %v",
				ctx.Command.Name, argsCount, ctx.Command.UsageText,
			)
		}
		return action(ctx)
	}
}
