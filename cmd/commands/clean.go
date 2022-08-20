package commands

import (
	"github.com/jieggii/dcfg/cmd/config"
	"github.com/jieggii/dcfg/cmd/log"
	"github.com/urfave/cli/v2"
	"os"
	p "path"
)

func Clean(ctx *cli.Context) error {
	cfgPath := ctx.String("config")
	dryRun := ctx.Bool("dry")

	cfg := config.ReadConfig(cfgPath)
	if dryRun {
		log.Info("* Dry run *")
	}
	log.Info("Cleaning context directory '%v':", cfg.Context)
	files, err := os.ReadDir(cfg.Context)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		name := file.Name()
		path := p.Join(cfg.Context, name)
		if path == cfgPath { // never remove dcfg config file
			log.Info("(~) Ignoring dcfg config file '%v'.", cfgPath)
			continue
		}
		if !cfg.PathIsDestination(name) {
			if !cfg.PathIsPinned(name) {
				var err error
				if dryRun {
					err = nil
				} else {
					err = os.RemoveAll(path)
				}
				if err == nil {
					log.Info("(-) Removing '%v'.", path)
				} else {
					log.Info("(!) Could not remove '%v' (%v).", path, err)
				}
			} else {
				log.Info("(~) Ignoring pinned '%v'.", path)
			}
		}
	}
	return nil
}
