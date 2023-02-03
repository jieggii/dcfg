package app

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

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
