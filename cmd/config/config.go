package config

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
)

type Bindings struct {
	// using own implementation of map because golang map behaves strangerly:
	// sometimes it messess up order of the keys, what leads to incorrect behaviour
	// of crossing bindings
	LongesetSourceLength int // this is needed for pretty-printing
	Sources              []string
	Destinations         []string
}

func (bindings *Bindings) AddBinding(source string, destination string) {
	sourceLen := len(source)
	if sourceLen > bindings.LongesetSourceLength {
		bindings.LongesetSourceLength = sourceLen
	}
	bindings.Sources = append(bindings.Sources, source)
	bindings.Destinations = append(bindings.Destinations, destination)
}

type Additions struct {
	LongestPathLenght int
	Paths             []string
}

func (additions *Additions) AddAddition(path string) {
	pathLen := len(path)
	if pathLen > additions.LongestPathLenght {
		additions.LongestPathLenght = pathLen
	}
	additions.Paths = append(additions.Paths, path)
}

type Config struct {
	Bindings  Bindings
	Additions Additions
}

type ParserError struct {
	LineNumber int
	Message    string
}

func unifyPath(path string, indicateDir bool) string {
	if strings.Contains(path, "~") { // expand user home dir (~)
		currentUser, err := user.Current()
		if err != nil {
			panic(err) // todo
		}
		path = strings.ReplaceAll(path, "~", currentUser.HomeDir)
	}
	if !strings.HasPrefix(path, "/") && !strings.HasPrefix(path, "./") {
		// this is just for pleasant look...
		// dir -> ./dir2
		path = "./" + path
	}
	if indicateDir {
		if !strings.HasSuffix(path, "/") {
			path = path + "/" // /dir -> /dir/
		}
	}
	return path
}

func unifyBindingPaths(path1 string, path2 string) (string, string) {
	if strings.HasSuffix(path1, "/") || strings.HasSuffix(path2, "/") {
		return unifyPath(path1, true), unifyPath(path2, true)
	} else {
		return unifyPath(path1, false), unifyPath(path2, false)
	}
}

func CreateConfig(path string) error {
	content := `# This is an example dcfg config file.
# More information about dcfg config files can be found here: https://github.com/jieggii/dcfg.

# "bind" directive - bind absolute path to relative one.
# Syntax: bind [absolute path] [relative path]
bind ~/ ./home/  # directories and files from $HOME will be copied to ./home/
bind / ./root/   # directories and files from / will be copied to ./root/

# "add" directive - copy directories and files to the current directory respecting bindings.
# Syntax: add [absolute path]
# add ~/.config/i3   # ~/.config/i3    will be copied to ./home/.config/i3
# add ~/.Xresources  # ~/.Xresourecs   will be copied to ./home/.Xresources
# add /etc/hostname  # /etc/hostname   will be copied to ./root/etc/hostname			   
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
	var config Config
	config.Bindings.LongesetSourceLength = 0
	config.Additions.LongestPathLenght = 0

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lineNumber := i + 1
		line = strings.Trim(line, " ")
		if line != "" && !strings.HasPrefix(line, "#") { // line is not empty and not full line comment
			valuable := strings.Split(line, "#")[0] // part of line before comment (e.g. "hello # world" -> "hello")
			valuable = strings.Trim(valuable, " ")  // removing extra spaces if they are present

			tokens := strings.Split(valuable, " ") // tokens of the valuable line
			directive := tokens[0]                 // first token = directive
			argsCount := len(tokens) - 1

			switch directive {
			case "bind": // bind [source] [dest]
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
				source, destination := unifyBindingPaths(tokens[1], tokens[2])
				config.Bindings.AddBinding(source, destination)
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
				path := unifyPath(tokens[1], false)
				config.Additions.AddAddition(path)
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
		fmt.Println("Error: there are some errors in the dcfg config file:")
		for _, parserError := range parserErrors {
			fmt.Printf("%v line %v: %v.\n", path, parserError.LineNumber, parserError.Message)
		}
		os.Exit(3)
	}
	return config, nil
}
