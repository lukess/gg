package cmds

import (
	"fmt"
	"ini"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	envp     = "/.gg/env"
	template = `
[core]
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
		fmt.Println("generating env file")
		os.Mkdir(".gg", 0755)
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

// from https://github.com/golang/go/blob/master/src/cmd/go/get.go
func runErr(cmd *exec.Cmd) error {
	out, err := cmd.CombinedOutput()
	if err != nil {
		if len(out) == 0 {
			return err
		}
		return fmt.Errorf("%s\n%v", out, err)
	}
	fmt.Println(string(out[:]))
	return nil
}

func Get() {
	wd, err := env_check()
	if err != nil {
		fmt.Println("env file does not exist")
	} else {
		ini.Init(wd + envp)
		libs := ini.GetAll("lib")
		for _, url := range libs {
			// url: lib[0], tag/branch: lib[1]
			lib := strings.Split(url, "=")
			if len(lib) != 2 {
				continue
			}
			// parse repo url
			repo_url := strings.TrimPrefix(strings.TrimPrefix(lib[0], "http://"), "https://")
			dir, file := filepath.Split(repo_url)
			filename := strings.TrimSuffix(file, filepath.Ext(file))
			// get the first element from GOPATH
			paths := strings.Split(ini.Get("core", "GOPATH"), ":")
			if len(paths) < 1 {
				fmt.Println("cannot read GOPATH")
				continue
			}
			repo_path := fmt.Sprintf("%s/src/%s", paths[0], dir)
			os.MkdirAll(repo_path, 0755)
			if strings.Contains(file, ".git") {
				runErr(exec.Command("git", "clone", "-b", lib[1], lib[0], repo_path+filename))
			} else {
				fmt.Printf("%s not a git repo\n", repo_url)
			}
		}
	}
}
