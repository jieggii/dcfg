package main

import (
	"github.com/jieggii/dcfg/internal/app"
	"github.com/jieggii/dcfg/internal/output"
	"os"
)

func main() {
	if err := app.NewApp().Run(os.Args); err != nil {
		output.Error.Printf("%v.\n", err)
		os.Exit(1)
	}
}
