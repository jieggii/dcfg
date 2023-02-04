package main

import (
	"github.com/jieggii/dcfg/internal/app"
	"github.com/jieggii/dcfg/internal/output"
	"os"
	"os/signal"
)

const interruptSignalExitCode = 130
const genericErrorExitCode = 1

func main() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	go func() {
		for range signalChannel {
			output.Stdout.Println()
			output.Error.Println("interrupt signal received")
			os.Exit(interruptSignalExitCode)
		}
	}()

	if err := app.NewApp().Run(os.Args); err != nil {
		output.Error.Println(err)
		os.Exit(genericErrorExitCode)
	}
}
