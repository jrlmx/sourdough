package main

import (
	"os"
	"os/exec"
)

func runCommand(cname string, cargs ...string) error {
	cmd := exec.Command(cname, cargs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
