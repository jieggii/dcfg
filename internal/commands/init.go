package commands

import (
	"fmt"
	"github.com/jieggii/dcfg/internal/config"
	"github.com/urfave/cli/v2"
	"os"
)

func Init(ctx *cli.Context) error {
	cfgPath := ctx.String("config")
	if _, err := os.Stat(cfgPath); err == nil { // check if file cfgPath already exists
		return fmt.Errorf("'%v' already exists", cfgPath)
	}

	cfg := config.NewConfig()
	if err := cfg.DumpToFile(cfgPath); err != nil {
		return err
	}

	return nil
}
