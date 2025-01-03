package main

import (
	"fmt"
	"os/exec"
)

func handlePHPDeps(cfg *config) error {
	fmt.Println("Installing composer dependencies...")

	if err := exec.Command("which", "composer").Run(); err != nil {
		return fmt.Errorf("composer is not installed")
	}

	if err := exec.Command("composer", "config", "repositories.flux-pro", "composer", "https://composer.fluxui.dev").Run(); err != nil {
		return fmt.Errorf("failed to add flux ui repository: %w", err)
	}

	deps := cfg.options.Packages.PHP

	if err := runCommand("composer", append([]string{"require"}, deps...)...); err != nil {
		return fmt.Errorf("failed to install composer dependencies: %w", err)
	}

	artisanCmds := cfg.options.Artisan.Commands

	for _, cmd := range artisanCmds {
		if err := runCommand("php", append([]string{"artisan"}, cmd)...); err != nil {
			return fmt.Errorf("failed to run %s: %v", cmd, err)
		}
	}

	return nil
}
