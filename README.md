About
=====
GG is a tool to create and manage Golang projects.

Feature:
* create isolated Golang environment
* create makefile template and directories
* get $GOPATH

gg Commands:

```bash
// add .gg/env file
$ gg setenv

// shell with env environment
$ gg env

// exec command in gg env
$ gg env make | gg getpath ...

// git clone libs defined in the .gg/env
$ gg get

// echo $GOPATH
$ gg getpath

// echo $GOROOT
$ gg getroot

// make file template
$ gg make

// create directories - dist/src
$ gg mkdirs
```

Build
=====
Modify OSS in Makefile if there is no cross compile environment.

	$ make
	$ sudo ln -s ~/dev/gg/dist/{OS}/gg /usr/bin/gg

Example
=====
1.create project directory

	$ mkdir myproject
	$ cd myproject

2.create bin/src/pkg and makefile

	$ gg mkdirs
	$ gg make

3.create gg config in .gg/env, and add a dependency

	$ gg setenv
	$ cat .gg/env
	$ echo "https://github.com/golang/example.git=master" >> .gg/env
	$ gg env
    #########################
    ## download dependency ##
    #########################
    $ gg get

4.create src/main.go

	gg:{myproject} {user}$ vim src/main.go

```go
package main

import (
	"fmt"
	"github.com/golang/example/stringutil"
)

func main() {
	fmt.Println(stringutil.Reverse("12345"))
}
```

5.build

	gg:{myproject} {user}$ make

6.run

	gg:{myproject} {user}$ dist/darwin_amd64/myproject
	54321

Misc
=====
[GOROOT environment variable is the path of Go binary distributions, and GOPATH specifies the location of your workspace.](https://golang.org/doc/code.html)

