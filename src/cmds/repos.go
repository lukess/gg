package cmds

import (
	"fmt"
	"ini"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// from https://github.com/golang/go/blob/master/src/cmd/go/get.go
func runErr(cmd *exec.Cmd, dir string) error {
	if dir != "" {
		cmd.Dir = dir
	}
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

/*
return {repo_url}, {repo_command}, {repo_target}
*/
func parseRepoMeta(repo_meta string) (string, string, string) {
	meta := strings.Split(repo_meta, ",")
	return meta[0], meta[1], meta[2]
}

/*
return
*/
func getRepoDir(repo_url string, createIfNotExist bool) string {
	// separate repo_url
	repo_plurl := strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(repo_url, "http://"), "https://"), "ssh://")
	dir, file := filepath.Split(repo_plurl)
	// support git only
	if !strings.Contains(file, ".git") {
		fmt.Printf("%s not a git repo\n", repo_plurl)
	}
	filename := strings.TrimSuffix(file, filepath.Ext(file))
	// get the first element from GOPATH
	paths := strings.Split(ini.Get("export", "GOPATH"), ":")
	if len(paths) < 1 {
		fmt.Errorf("cannot read GOPATH")
	}
	repo_path := fmt.Sprintf("%s/src/%s", paths[0], dir)
	if createIfNotExist {
		os.MkdirAll(repo_path, 0755)
	}
	return repo_path + filename
}

/*
lib section format:
{repo_id}={github repo url},{branch/tag/commit},{branch name/tag name/commit hash}
*/
func Get() {
	wd, err := env_check()
	if err != nil {
		fmt.Println("env file does not exist")
	} else {
		ini.Init(wd + envp)
		libs := ini.GetAll("lib")
		for _, repo := range libs {
			meta := strings.Split(repo, "=")
			if len(meta) != 2 {
				continue
			}
			// parse repo
			repo_url, repo_command, repo_target := parseRepoMeta(meta[1])
			local_repodir := getRepoDir(repo_url, true)
			// branch, tag or commit
			if repo_command == "branch" {
				runErr(exec.Command("git", "clone", "-b", repo_target, repo_url, local_repodir), "")
			} else if repo_command == "tag" || repo_command == "commit" {
				runErr(exec.Command("git", "clone", repo_url, local_repodir), "")
				runErr(exec.Command("git", "checkout", repo_target), local_repodir)
			} else {
				fmt.Printf("%s does not support\n", repo_command)
			}
		}
	}
}

func Git(args []string) {
	wd, err := env_check()
	if err != nil {
		fmt.Println("env file does not exist")
	} else {
		ini.Init(wd + envp)
		libs := ini.GetAll("lib")
		repo_id := args[0]
		for _, repo := range libs {
			meta := strings.Split(repo, "=")
			if len(meta) != 2 {
				continue
			}
			if repo_id == meta[0] {
				repo_url, _, _ := parseRepoMeta(meta[1])
				local_repodir := getRepoDir(repo_url, false)
				runErr(exec.Command(args[1], args[2:]...), local_repodir)
			}
		}
	}
}
