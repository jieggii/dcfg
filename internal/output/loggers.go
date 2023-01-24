package output

import (
	"log"
	"os"
)

var Info = log.New(os.Stdout, "", 0)
var Warning = log.New(os.Stdout, "warning: ", 0)
var Error = log.New(os.Stderr, "error: ", 0)
