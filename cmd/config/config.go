package config

import (
	"errors"
	"fmt"
	"github.com/jieggii/dcfg/cmd/log"
	"github.com/jieggii/dcfg/cmd/util"
	"os"
	p "path"
	"strconv"
	"strings"
)

type ParserError struct {
	LineNumber int
	Message    string
}

func CreateConfig(path string) error {
	content := `# This is an example dcfg config file.
# More information about dcfg config files can be found here: https://github.com/jieggii/dcfg.

# 'dest' directive - set directory where all config files will be placed.
# Syntax: dest [path].
dest ./  # ./ is default value

# Bindings (order makes sense):
# 'bind' directive - bind absolute path to a local one.
# Syntax: bind [absolute path] [local path].
bind ~ home/  # directories and files from $HOME will be copied to ./home/
bind / root/  # directories and files from / will be copied to ./root/

# Additions:
# 'add' directive - copy directories and files to the destination directory respecting bindings.
# Syntax: add [absolute path]
# add ~/.config/i3   # ~/.config/i3  will be copied to ./home/.config/i3
# add ~/.Xresources  # ~/.Xresources will be copied to ./home/.Xresources
# add /etc/hostname  # /etc/hostname will be copied to ./root/etc/hostname

`
	if _, err := os.Stat(path); err == nil {
		return errors.New("file already exists")
	} else if errors.Is(err, os.ErrNotExist) {
		return os.WriteFile(path, []byte(content), 0666) // todo: set proper perm
	} else {
		return err
	}
}

func parseConfig(content string) (*Config, []ParserError) {
	var parserErrors []ParserError
	config := newConfig()

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lineNumber := i + 1
		line = strings.Trim(line, " ")
		if line != "" && !strings.HasPrefix(line, "#") { // line is not empty and not full line comment
			valuable := strings.Split(line, "#")[0] // part of line before comment (e.g. "hello # world" -> "hello")
			valuable = strings.Trim(valuable, " ")  // removing extra spaces if they are present

			tokens := strings.Split(valuable, " ") // tokens of the valuable line
			directive := tokens[0]                 // first token - directive

			args := tokens[1:] // other tokens - arguments
			argsCount := len(args)

			switch directive {
			case "dest": // dest [path]
				if argsCount != 1 {
					parserErrors = append(
						parserErrors,
						ParserError{
							lineNumber,
							"'dest' directive requires exactly one argument, got " + strconv.Itoa(argsCount),
						},
					)
					break
				}
				if config.Bindings.DestinationPrefixWasSet {
					parserErrors = append(
						parserErrors,
						ParserError{
							lineNumber,
							"'destination' directive may be used only once",
						},
					)
					break
				}
				path := util.CompilePath(args[0])
				if p.IsAbs(path) {
					parserErrors = append(
						parserErrors,
						ParserError{
							lineNumber,
							"the only argument of 'destination' directive (path) must be a local path",
						},
					)
					break
				}
				config.Bindings.DestinationPrefix = path
				config.Bindings.DestinationPrefixWasSet = true

			case "bind": // bind [source] [dest]
				if argsCount != 2 {
					parserErrors = append(
						parserErrors,
						ParserError{
							lineNumber,
							"'bind' directive requires exactly two arguments, got " + strconv.Itoa(argsCount),
						},
					)
					break
				}
				source := util.CompilePath(args[0])
				if !p.IsAbs(source) {
					parserErrors = append(
						parserErrors,
						ParserError{
							lineNumber,
							"the first argument of 'bind' directive (source path) must be an absolute path",
						},
					)
					break
				}
				destination := args[1]
				if p.IsAbs(destination) {
					parserErrors = append(
						parserErrors,
						ParserError{
							lineNumber,
							"the second argument of 'bind' directive (destination path) must be a local path",
						},
					)
				}
				config.Bindings.appendBinding(source, destination)
			case "add": // bind [path]
				if argsCount != 1 {
					parserErrors = append(
						parserErrors,
						ParserError{
							lineNumber,
							"'add' directive requires exactly one argument, got " + strconv.Itoa(argsCount),
						})
					break
				}
				path := util.CompilePath(args[0])
				if !p.IsAbs(path) {
					parserErrors = append(
						parserErrors,
						ParserError{
							lineNumber,
							"the only argument of 'add' directive (path) must be an absolute path",
						},
					)
					break
				}
				config.Additions.appendAddition(path)
			default:
				parserErrors = append(
					parserErrors,
					ParserError{
						lineNumber,
						fmt.Sprintf("unknown directive '%v'", directive),
					},
				)
			}
		}
	}
	return config, parserErrors
}

func ReadConfig(path string) *Config {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Error("Error: could not read dcfg config file %v: %v.", path, err)
		os.Exit(2)
	}
	config, parserErrors := parseConfig(string(bytes))
	if len(parserErrors) != 0 {
		for _, parserError := range parserErrors {
			log.Error("Error parsing %v (line %v): %v.", path, parserError.LineNumber, parserError.Message)
		}
		os.Exit(3)
	}
	return config
}
