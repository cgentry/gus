// Copyright 2014 Charles Gentry All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// gus is the Go User Service. It provides a simple interface for login, logout, authenticate
// and general user services. Is is flexible and configurable. The data can be stored in a number
// of different ways (database, flat files) and different configurations can be selected.
// The system has been carefully layered to segregate each part for easier expansion. You can
// Change the way the system is called (an http is the current interface), change where the
// data is stored, the way keys are encrypted, and you can limit the available options for
// user updates.
//
// GUS has a built-in configuration and bootstrap system that can be used as a command-line
// program to setup the system for the first time. It includes extensive help and interactive
// prompts to help you configure the system.
//
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/cgentry/gus/cli"
	"github.com/cgentry/gus/record/configure"
)

var configFileName string
var commands = []*cli.Command{
	cmdConfig,
	cmdCreateStore,
	cmdUser,
	cmdUserAdd,
	cmdService,
	helpStore,
	helpEncrypt,
}

var helpTemplate = `Usage:

          go command [arguments]

{{range .}}
    {{.Id | printf "%-11s"}} {{ .ShortHelp }}{{end}}

Use "gus help [command]" for more information about a cli.

`

// Usage will output a help text and then exit.
func Usage() {
	cli.Usage(helpTemplate, commands)
	os.Exit(0)
}

func main() {
	flag.Usage = Usage
	flag.Parse()
	log.SetFlags(0)

	args := flag.Args()
	if len(args) < 1 {
		Usage()
	}
	if args[0] == "help" {
		cli.Help(helpTemplate, "gus", args, commands)
		return
	}

	// Try and run the command given...
	for _, cmd := range commands {
		if cmd.Name == args[0] {
			cmd.Flag.Usage = func() { cmd.Usage() }
			if cmd.CustomFlags {
				cmd.Run(cmd, args[1:])
			} else {
				cmd.Flag.Parse(args[1:])
				args = cmd.Flag.Args()
				cmd.Run(cmd, args)
			}
			return
		}
	}
}

// GetConfigFileName is used to get the configuration filename and check the existance of the file. It will return
// errors if either the directory or the file doesn't exist.
func GetConfigFileName() (Filename string, DirExists error, FileExists error) {
	if configFileName == "" {
		configFileName = DefaultConfigFilename
	}
	_, DirExists = os.Stat(filepath.Dir(configFileName))
	_, FileExists = os.Stat(configFileName)
	Filename = configFileName
	return
}

// GetConfigFile will pull open up the configuration file and return a configuration object.
func GetConfigFile() (*configure.Configure, error) {
	var err error
	c := configure.New()

	fname, dirError, fileError := GetConfigFileName()
	if dirError != nil {
		return c, dirError
	}
	if fileError != nil {
		return c, fileError
	}
	fdata, err := ioutil.ReadFile(fname)
	if err == nil {
		err = json.Unmarshal(fdata, c)
	}

	return c, err
}

// SaveConfigFile will take the JSON string and write it to the output file
func SaveConfigFile(jsonString []byte) error {
	file, direrror, _ := GetConfigFileName()
	if direrror != nil {
		return direrror
	}
	return ioutil.WriteFile(file, jsonString, DefaultConfigPermissions)
}

// addCommonCommandFlags will add in flags that are system-wide.
func addCommonCommandFlags(cmd *cli.Command) {
	cmd.Flag.StringVar(&configFileName, "c", DefaultConfigFilename, "")
}

func runtimeFail(msg string, err error) {
	var rpt int
	var emsg string
	if err == nil {
		emsg = "(runtime error)"
	} else {
		emsg = err.Error()
	}
	if len(emsg) > len(msg) {
		rpt = len(emsg)
	} else {
		rpt = len(msg)
	}

	stars := strings.Repeat("*", rpt+4)
	fmt.Fprintf(os.Stderr, "%s\n* %-*s *\n* %-*s *\n%s\n\n", stars, rpt, msg, rpt, emsg, stars)
	fmt.Fprintln(os.Stderr, "STACK TRACE:")
	debug.PrintStack()
	fmt.Fprintln(os.Stderr, "\n")
	os.Exit(1)
}
