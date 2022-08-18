package commands

import (
	"github.com/jieggii/dcfg/cmd/config"
	"github.com/jieggii/dcfg/cmd/log"
	"github.com/jieggii/dcfg/cmd/util"
	"github.com/urfave/cli/v2"
	"os"
)

func Init(context *cli.Context) error {
	cfgPath := context.String("config")
	if context.Args().Len() != 0 {
		util.ThrowUsageError(context, "too many arguments")
	}
	if err := config.CreateConfig(cfgPath); err != nil {
		log.Error("Error: could not create dcfg config file `%v`: %v.", cfgPath, err)
		os.Exit(2)
	}
	log.Info("Created dcfg config file `%v` in the current working directory.", cfgPath)
	return nil
}
