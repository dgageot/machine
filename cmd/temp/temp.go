package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("echo", "$TOTO")
	cmd.Env = append(os.Environ(), "TOTO="+"DAVID")

	out, _ := cmd.Output()

	fmt.Println("[" + string(out) + "]")
}
