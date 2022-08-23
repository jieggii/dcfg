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
	dryRun := context.Bool("dry")
	util.CheckCommandArgsCount(context, 0)
	if dryRun {
		util.LogDryRun()
	}

	var err error
	if dryRun {
		err = nil
	} else {
		err = config.CreateConfig(cfgPath)
	}
	if err != nil {
		log.Error("Error: could not create dcfg config file `%v`: %v.", cfgPath, err)
		os.Exit(2)
	}
	log.Info("Created dcfg config file `%v` in the current working directory.", cfgPath)
	return nil
}
