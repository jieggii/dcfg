package commands

import (
	"fmt"
	"github.com/jieggii/dcfg/cmd/config"
	"github.com/jieggii/dcfg/cmd/log"
	cp "github.com/otiai10/copy"
	"github.com/urfave/cli/v2"
	"os"
	p "path"
	"strconv"
	"strings"
)

func Update(context *cli.Context) error {
	cfgPath := context.String("config")
	cfg := config.ReadConfig(cfgPath)
	if len(cfg.Additions.Paths) == 0 {
		log.Error("Error: nothing to do (there are no `add` directives in `%v`).", cfgPath)
		os.Exit(0)
	}
	log.Info(
		"Updating destination directory (%v) according to these bindings:",
		cfg.Bindings.DestinationPrefix,
	)
	for i, source := range cfg.Bindings.Sources {
		destination := cfg.Bindings.Destinations[i]
		log.Info("%v. %-"+strconv.Itoa(cfg.Bindings.LongestSourceLength)+"v : %v", i+1, source, destination)
	}
	fmt.Println("")
	for _, globalPath := range cfg.Additions.Paths {
		matched := false
		for i, source := range cfg.Bindings.Sources {
			destination := cfg.Bindings.Destinations[i]
			if strings.HasPrefix(globalPath, source) {
				matched = true
				localPath := strings.Replace(globalPath, source, destination, 1)
				localPath = p.Clean(localPath)
				localPath = p.Join(cfg.Bindings.DestinationPrefix, localPath)
				if err := cp.Copy(globalPath, localPath); err != nil {
					log.Info("Warning: could not copy %v to %v (%v)", globalPath, localPath, err)
				} else {
					log.Info("%-"+strconv.Itoa(cfg.Additions.LongestPathLength)+"v -> %v", globalPath, localPath)
				}
				break
			}
		}
		if !matched {
			log.Info("Warning: ignoring unmatched addition: %v", globalPath)
		}
	}
	log.Info("")
	log.Info("Removing non-destination dirs (as --remove option is used):")
	if context.Bool("remove") {
		files, err := os.ReadDir(cfg.Bindings.DestinationPrefix)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			name := file.Name()
			if name != cfgPath && name != ".git" {
				isDestination := false
				for _, destination := range cfg.Bindings.Destinations {
					if name == p.Base(destination) {
						isDestination = true
						break
					}
				}
				if !isDestination {
					log.Info("Remove %v", name)
				}
			}
		}
	}

	return nil
}
