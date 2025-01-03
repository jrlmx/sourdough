package main

import (
	"fmt"
	"os/exec"
)

func handleJSDeps(cfg *config) error {
	if err := exec.Command("which", "npm").Run(); err != nil {
		return fmt.Errorf("npm is not installed")
	}

	deps := cfg.options.Packages.JS

	if err := runCommand("npm", append([]string{"install", "-D"}, deps...)...); err != nil {
		return fmt.Errorf("failed to install node dependencies: %w", err)
	}

	return nil
}
