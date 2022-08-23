package commands

import (
	"github.com/jieggii/dcfg/cmd/config"
	"github.com/jieggii/dcfg/cmd/log"
	"github.com/jieggii/dcfg/cmd/util"
	"github.com/urfave/cli/v2"
	"os"
	p "path"
)

func Clean(ctx *cli.Context) error {
	cfgPath := ctx.String("config")
	dryRun := ctx.Bool("dry")
	cfg := config.ReadConfig(cfgPath)
	if dryRun {
		util.LogDryRun()
	}

	log.Info("Cleaning context directory '%v'.", cfg.Context)
	files, err := os.ReadDir(cfg.Context)
	if err != nil {
		log.Error("Error reading context directory: %v.", err)
		os.Exit(3)
	}
	for _, file := range files {
		name := file.Name()               // name of the file or dir inside context dir
		path := p.Join(cfg.Context, name) // path to the file or dir relative to workdir
		if path == cfgPath {
			log.Info("(~) Ignoring dcfg config file '%v'.", path)
			continue
		}
		if name == ".git" {
			log.Info("(~) Ignoring git directory '%v'.", path)
			continue
		}
		if cfg.PathIsDestination(name) {
			log.Info("(~) Ignoring destination '%v'.", path)
			continue
		}
		if cfg.PathIsPin(name) {
			log.Info("(~) Ignoring pinned '%v'.", path)
			continue
		}
		var err error
		if dryRun {
			err = nil
		} else {
			err = os.RemoveAll(path)
		}
		if err != nil {
			log.Info("(!) Could not remove '%v' (%v).", path, err)
		} else {
			log.Info("(-) Removing '%v'.", path)
		}
	}
	return nil
}
