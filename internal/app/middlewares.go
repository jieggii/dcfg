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
				"%v command takes from %v to %v arguments, got %v\nusage: %v",
				ctx.Command.Name, minArgsCount, maxArgsCount, argsCount, ctx.Command.UsageText,
			)
		}
	}
}

func moreThanArgsCountMiddleware(minimalArgsCount int, action cli.ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		argsCount := ctx.Args().Len()
		if argsCount < minimalArgsCount {
			return fmt.Errorf(
				"%v command takes %v and more arguments, got %v\nusage: %v",
				ctx.Command.Name, minimalArgsCount, argsCount, ctx.Command.UsageText,
			)
		}
		return action(ctx)
	}
}

func explicitArgsCountMiddleware(expectedArgsCount int, action cli.ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		argsCount := ctx.Args().Len()
		if argsCount == expectedArgsCount {
			return action(ctx)
		} else {
			if expectedArgsCount == 0 {
				return fmt.Errorf(
					"%v command takes no arguments, got %v\nusage: %v",
					ctx.Command.Name, argsCount, ctx.Command.UsageText,
				)
			} else {
				return fmt.Errorf(
					"%v command takes exactly %v argument(s), got %v\nusage: %v",
					ctx.Command.Name, expectedArgsCount, argsCount, ctx.Command.UsageText,
				)
			}
		}
	}
}
