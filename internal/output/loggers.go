package output

import (
	"log"
	"os"
)

var Stdout = log.New(os.Stdout, "", 0)

var Info = log.New(os.Stdout, "[*] ", 0)
var Success = log.New(os.Stdout, "[+] ", 0)

var Error = log.New(os.Stderr, "error: ", 0)
