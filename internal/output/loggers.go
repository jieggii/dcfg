package output

import (
	"log"
	"os"
)

var Stdout = log.New(os.Stdout, "", 0)

var Plus = log.New(os.Stdout, "(+) ", 0)
var Minus = log.New(os.Stdout, "(-) ", 0)

var Warning = log.New(os.Stdout, "(!) ", 0)

// Error logger is used only on fatal errors. Otherwise, Warning logger is used.
var Error = log.New(os.Stderr, "error: ", 0)
