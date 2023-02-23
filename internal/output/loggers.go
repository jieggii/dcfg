package output

import (
	"log"
	"os"
)

var Stdout = log.New(os.Stdout, "", 0)

var Plus = log.New(os.Stdout, "(+) ", 0)
var Minus = log.New(os.Stdout, "(-) ", 0)
var Verbose = log.New(os.Stdout, "(v) ", 0)

var Warning = log.New(os.Stdout, "(!) ", 0)

var Error = log.New(os.Stderr, "error: ", 0)
