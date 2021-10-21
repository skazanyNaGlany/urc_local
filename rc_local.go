package main

import (
	"log"
	"os/exec"
)

const Arg0 = "sh"
const Arg1 = "-x"
const Arg2 = "rc_local.sh"

func main() {
	output, err := exec.Command(Arg0, Arg1, Arg2).Output()

	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(output))
}
