package commands

import (
	"fmt"
	"github.com/jieggii/dcfg/cmd/config"
	"github.com/jieggii/dcfg/cmd/log"
	"github.com/jieggii/dcfg/cmd/util"
	cp "github.com/otiai10/copy"
	"github.com/urfave/cli/v2"
	"os"
	p "path"
	"strconv"
	"strings"
)

func Add(context *cli.Context) error {
	cfgPath := context.String("config")
	dryRun := context.Bool("dry")
	cfg := config.ReadConfig(cfgPath)
	if dryRun {
		util.LogDryRun()
	}

	if len(cfg.Additions.Paths) == 0 {
		log.Error("Error: nothing to do (there are no `add` directives in `%v`).", cfgPath)
		os.Exit(0)
	}
	if len(cfg.Bindings.Sources) == 0 {
		log.Info("Copying additions to the context directory '%v'.", cfg.Context)

	} else {
		log.Info(
			"Copying additions to the context directory '%v' according to these bindings:",
			cfg.Context,
		)
	}
	for i, source := range cfg.Bindings.Sources {
		destination := cfg.Bindings.Destinations[i]
		log.Info("%v. %-"+strconv.Itoa(cfg.Bindings.LongestSourceLength)+"v : %v", i+1, source, destination)
	}
	if len(cfg.Bindings.Sources) != 0 {
		fmt.Println("")
	}
	for _, globalPath := range cfg.Additions.Paths {
		matched := false
		for i, source := range cfg.Bindings.Sources {
			destination := cfg.Bindings.Destinations[i]
			if strings.HasPrefix(globalPath, source) {
				matched = true
				localPath := strings.Replace(globalPath, source, destination, 1)
				localPath = p.Clean(localPath)
				localPath = p.Join(cfg.Context, localPath)
				var err error
				if dryRun {
					err = nil
				} else {
					err = cp.Copy(globalPath, localPath)
				}
				if err != nil {
					log.Info("Warning: could not copy '%v' to '%v' (%v)", globalPath, localPath, err)
				} else {
					log.Info("(+) %-"+strconv.Itoa(cfg.Additions.LongestPathLength)+"v -> %v", globalPath, localPath)
				}
				break
			}
		}
		if !matched {
			log.Info("Warning: ignoring unmatched addition: %v", globalPath)
		}
	}
	return nil
}
