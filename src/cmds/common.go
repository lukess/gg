package cmds

import (
	"log"
	"os/exec"
)

var ()

func Shell(cmd string) []byte {
	p, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Fatal(err)
	}
	return p
}
