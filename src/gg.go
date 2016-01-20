package main

import (
	"cmds"
	"fmt"
	"os"
)

var ()

func stringInSlice(compare []string, list []string) bool {
	for _, b := range list {
		for _, a := range compare {
			if b == a {
				return true
			}
		}
	}
	return false
}

func main() {
	help := `gg [OPTION]... [COMMANDS] args...

OPTIONS:
  -h print current help

COMMANDS:
  setenv
    create .gg/env if the file does not exist

  env
    interactive shell with env variables from .gg/env

  getpath
    echo $GOPATH

  getroot
    echo $GOROOT

  make
    generate a simple Makefile template within current directory

  mkdirs
    create bin, pkg and src directories
`
	if stringInSlice([]string{"-h", "--help"}, os.Args) {
		fmt.Println(help)
		return
	}

	if len(os.Args) > 1 {
		command := os.Args[1]
		switch {
		case command == "setenv":
			cmds.Setenv()
		case command == "env":
			cmds.Env(os.Args[2:])
		case command == "getpath":
			fmt.Println(os.Getenv("GOPATH"))
		case command == "getroot":
			fmt.Println(os.Getenv("GOROOT"))
		case command == "get":
			cmds.Get()
		case command == "git":
			cmds.Git(os.Args[2:])
		case command == "make":
			cmds.Make()
		case command == "mkdirs":
			cmds.Mkdirs()
		default:
			fmt.Println("gg command not found")
		}
	}
}
