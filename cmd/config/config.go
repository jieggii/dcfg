package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Bindings  map[string]string
	Additions []string
}

type ParserError struct {
	LineNumber int
	Message    string
}

func CreateConfig(path string) error {
	content := `# This is an example dcfg config file.
# More information about dcfg config files can be found here: https://github.com/jieggii/dcfg.

# "bind" directive - bind absolute path to local one.
# Syntax: bind [absolute path] [relative path]
bind ~ ./home/  # directories and files from $HOME will be copied to ./home/
bind / ./root/  # directories and files from / will be copied to ./root/

# "add" directive - copy directories and files to the current directory respecting bindings.
# Syntax: add [absolute path]
# add ~/.config/i3     # will be copied to ./home/.config/i3
# add ~/.config/picom  # will be copied to ./home/.config/picom
# add ~/.Xresources    # will be copied to ./home/.Xresources
`
	if _, err := os.Stat(path); err == nil {
		return errors.New("file already exists")
	} else if errors.Is(err, os.ErrNotExist) {
		return os.WriteFile(path, []byte(content), 0755)
	} else {
		return err
	}
}

func parseConfig(content string) (*Config, []ParserError) {
	var config Config
	config.Bindings = make(map[string]string) // init bindings map

	var parserErrors []ParserError

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lineNumber := i + 1
		line = strings.Trim(line, " ")
		if line != "" && !strings.HasPrefix(line, "#") { // line is not empty and not full line comment
			valuable := strings.Split(line, "#")[0] // part of line before comment (e.g. "hello # world" -> "hello")
			valuable = strings.Trim(valuable, " ")  // removing extra spaces if they are present

			tokens := strings.Split(valuable, " ") // tokens of the valuable line
			directive := tokens[0]                 // first token - directive
			argsCount := len(tokens) - 1

			switch directive {
			case "bind": // bind [path1] [path2]
				if argsCount != 2 {
					parserErrors = append(
						parserErrors,
						ParserError{
							lineNumber,
							"`bind` directive requires exactly two arguments, got " + strconv.Itoa(argsCount),
						},
					)
					break
				}
				path1 := tokens[1]
				path2 := tokens[2]
				config.Bindings[path1] = path2
			case "add": // bind [path]
				if argsCount != 1 {
					parserErrors = append(
						parserErrors,
						ParserError{
							lineNumber,
							"`add` directive requires exactly one argument, got " + strconv.Itoa(argsCount),
						})
					break
				}
				path1 := tokens[1]
				config.Additions = append(config.Additions, path1)
			default:
				parserErrors = append(
					parserErrors,
					ParserError{
						lineNumber,
						fmt.Sprintf("unknown directive `%v`", directive),
					},
				)
			}
		}
	}
	return &config, nil
}

func ReadConfig(path string) (*Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config, parserErrors := parseConfig(string(bytes))
	if len(parserErrors) != 0 {
		fmt.Println("Error: got some errors while parsing dcfg config file:")
		for _, parserError := range parserErrors {
			fmt.Printf("%v line %v: %v.\n", path, parserError.LineNumber, parserError.Message)
		}
		os.Exit(2)
	}
	return config, nil
}
