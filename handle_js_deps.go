package main

import (
	"fmt"
	"os/exec"
)

func handleJSDeps(cfg *config) error {
	if err := exec.Command("which", "npm").Run(); err != nil {
		return fmt.Errorf("npm is not installed")
	}

	prod := cfg.opts.JS.Prod

	if err := runCommand("npm", append([]string{"install", "-D"}, prod...)...); err != nil {
		return fmt.Errorf("failed to install node dependencies: %w", err)
	}

	dev := cfg.opts.JS.Dev

	if err := runCommand("npm", append([]string{"install", "-D"}, dev...)...); err != nil {
		return fmt.Errorf("failed to install node dependencies: %w", err)
	}

	return nil
}
