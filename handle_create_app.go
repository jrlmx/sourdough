package main

import (
	"fmt"
	"os"
	"os/exec"
)

func handleCreateApp(cfg *config) error {
	if err := exec.Command("which", "laravel").Run(); err != nil {
		return fmt.Errorf(("laravel installer is not installed"))
	}

	if _, err := os.Stat(cfg.projectName); err == nil {
		return fmt.Errorf("directory already exists: %s", cfg.projectDir)
	}

	fmt.Printf("Creating new Laravel project: %s\n", cfg.projectName)

	cmd := exec.Command("laravel", "new", cfg.projectName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create Laravel project: %w", err)
	}

	if err := os.Chdir(cfg.projectDir); err != nil {
		return fmt.Errorf("failed to change to project directory: %w", err)
	}

	return nil
}
