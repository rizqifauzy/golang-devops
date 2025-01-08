package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("/bin/bash", "-c", "df")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
	} else {
		fmt.Printf("Command Output:\n%s\n", string(output))
	}
}
