package log

import (
	"fmt"
	"io"
	"os"
)

func writeln(writer io.Writer, message string, a ...any) {
	if _, err := fmt.Fprintf(writer, message+"\n", a...); err != nil {
		panic(err)
	}
}

func Error(message string, a ...any) {
	writeln(os.Stderr, message, a...)
}

func Info(message string, a ...any) {
	writeln(os.Stdout, message, a...)
}
