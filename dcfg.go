package dcfg

import (
	"github.com/jieggii/dcfg/internal"
	"github.com/jieggii/dcfg/internal/app"
	"github.com/jieggii/dcfg/internal/output"
	"os"
	"os/signal"
)

func main() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	go func() {
		for range signalChannel {
			output.Stdout.Println()
			output.Stdout.Println("interrupt signal received")
			os.Exit(internal.InterruptSignalExitCode)
		}
	}()

	if err := app.NewApp().Run(os.Args); err != nil {
		output.Error.Println(err)
		os.Exit(internal.GenericErrorExitCode)
	}
}
