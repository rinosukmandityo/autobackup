package helper

import (
	"log"
	"os/exec"
)

func ExecCommand(commands []string) {
	if e := exec.Command("cmd", commands...).Run(); e != nil {
		log.Println(e.Error())
	}
	log.Println("exec command success")
}
