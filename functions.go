package main

import (
	"fmt"
	"os"
	"os/exec"
)

func run(cname string, cargs ...string) error {
	cmd := exec.Command(cname, cargs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func check(cname string) error {
	if err := exec.Command("which", "npm").Run(); err != nil {
		return fmt.Errorf("%s is not installed", cname)
	}

	return nil
}
