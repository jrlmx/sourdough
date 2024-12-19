package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func handleComposerDeps(cfg *config) error {
	if err := exec.Command("which", "composer").Run(); err != nil {
		return fmt.Errorf("composer is not installed")
	}

	if err := exec.Command("composer", "config", "repositories.flux-pro", "composer", "https://composer.fluxui.dev").Run(); err != nil {
		return fmt.Errorf("failed to add flux ui repository: %w", err)
	}

	deps := cfg.deps.Composer

	cmd := exec.Command("composer", append([]string{"require"}, deps...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install composer dependencies: %w", err)
	}

	artisanCmds := []string{
		"folio:install",
		"volt:install",
	}

	for _, artisanCmd := range artisanCmds {
		cmd := exec.Command("php", "artisan", artisanCmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Failed to run %s: %v", artisanCmd, err)
		}
	}

	return nil
}
