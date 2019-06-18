package test

import (
	"fmt"
	"os"
	"os/exec"
)

// ExecShell executes a shell command
func ExecShell(cmd string) error {
	fmt.Println(cmd)
	builtCmd := exec.Command("sh", "-c", cmd)
	builtCmd.Stderr = os.Stderr
	builtCmd.Stdout = os.Stdout
	return builtCmd.Run()
}
