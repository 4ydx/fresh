/*
Fresh is a command line tool that builds and (re)starts your web application everytime you save a go or template file.

If the web framework you are using supports the Fresh runner, it will show build errors on your browser.

It currently works with Traffic (https://github.com/pilu/traffic), Martini (https://github.com/codegangsta/martini) and gocraft/web (https://github.com/gocraft/web).

Fresh will watch for file events, and every time you create/modifiy/delete a file it will build and restart the application.
If `go build` returns an error, it will logs it in the tmp folder.

Traffic (https://github.com/pilu/traffic) already has a middleware that shows the content of that file if it is present. This middleware is automatically added if you run a Traffic web app in dev mode with Fresh.
*/
package main

import (
	"flag"
	"fmt"
	"github.com/4ydx/fresh/runner"
	"os"
)

func main() {
	arguments, config := trimArguments(os.Args)

	os.Args = config // removed all other values so that flag parsing wont complain
	configPath := flag.String("c", "", "config file path")
	flag.Parse()

	if *configPath != "" {
		if _, err := os.Stat(*configPath); err != nil {
			fmt.Printf("Can't find config file `%s`\n", *configPath)
			os.Exit(1)
		} else {
			os.Setenv("RUNNER_CONFIG_PATH", *configPath)
		}
	}

	runner.Start(arguments)
}

// trimArguments returns arguments to pass to the binary being rebuilt
// as well as values to be directly consumed by flag parsing
func trimArguments(arguments []string) ([]string, []string) {
	program := arguments[0]
	arguments = arguments[1:]
	removeAt := -1
	for i, arg := range arguments {
		if arg == "c" {
			if removeAt != -1 {
				panic("Too many config arguments specified")
			}
			removeAt = i
		}
	}

	// rebuild a sane list of arguments for flag to parse
	config := make([]string, 0)
	config = append(config, program)

	// get our list of arguments to pass to our binary that we restart
	trimmed := make([]string, 0)
	if removeAt == -1 {
		trimmed = arguments
	} else {
		for i, arg := range arguments {
			if i != removeAt && i != removeAt+1 {
				trimmed = append(trimmed, arg)
			} else {
				config = append(config, arg)
			}
		}
	}
	return trimmed, config
}
