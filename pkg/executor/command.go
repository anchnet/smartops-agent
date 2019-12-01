package executor

import (
	"errors"
	"fmt"
	"os/exec"
)

func RunCommand(cmdStr string) (string, error) {
	cmd := exec.Command("bash", "-c", cmdStr)
	out, err := cmd.Output()
	if err != nil {
		e := err.(*exec.ExitError)
		return "", errors.New(string(e.Stderr))
	}
	return fmt.Sprint(string(out)), nil
}
