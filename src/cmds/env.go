package cmds

import (
	"fmt"
	"ini"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

var (
	envp     = "/gg"
	template = `
[export]
PATH=/usr/bin:/bin:%s/bin
PS1=\e[0;32mgg\e[m:\w \u\$
GOROOT=%s
GOPATH=%s
HOME=%s
TERM=%s
[lib]
`
)

func env_check() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(wd + envp); err != nil {
		return wd, err
	}
	return wd, nil
}

func Setenv() {
	wd, err := env_check()
	if err != nil {
		fmt.Println(err)
		fmt.Println("generating gg file")
		root := os.Getenv("GOROOT")
		if root == "" {
			if runtime.GOOS != "windows" {
				// default path
				root = "/usr/local/go"
			} else {
				root = "C:\\go"
			}
		}
		env := fmt.Sprintf(template, os.Getenv("GOROOT"), root, wd, os.Getenv("HOME"), os.Getenv("TERM"))
		if err := ioutil.WriteFile(wd+envp, []byte(env), 0644); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("env file exists")
	}
}

func Env(args []string) {
	wd, err := env_check()
	if err != nil {
		fmt.Println("env file does not exist")
	} else {
		ini.Init(wd + envp)
		env := []string{}
		for _, e := range ini.GetAll("export") {
			if pairs := strings.Split(e, "="); len(pairs) == 2 {
				if pairs[0] == "PS1" {
					env = append(env, fmt.Sprintf("%s=%s ", pairs[0], pairs[1]))
				} else {
					env = append(env, fmt.Sprintf("%s=%s", pairs[0], pairs[1]))
				}
			}
		}
		pa := os.ProcAttr{
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
			Dir:   wd,
			Env:   env,
		}
		var proc *os.Process
		if runtime.GOOS != "windows" {
			if len(args) == 0 {
				proc, err = os.StartProcess(os.Getenv("SHELL"), []string{"sh"}, &pa)
			} else {
				proc, err = os.StartProcess(os.Getenv("SHELL"), []string{"sh", "-c", strings.Join(args, " ")}, &pa)
			}
		}

		if err != nil {
			panic(err)
		}
		_, err = proc.Wait()
		if err != nil {
			panic(err)
		}
	}
}
