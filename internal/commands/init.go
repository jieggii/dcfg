package commands

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func Init(ctx *cli.Context) error {
	config := ctx.String("config")
	fmt.Println(config)
	return nil
}
