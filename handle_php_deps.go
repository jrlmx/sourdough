package main

import (
	"fmt"
	"strings"
)

func handlePHPDeps(p *project) error {
	fmt.Println("Installing composer dependencies...")
	for _, repo := range p.config.Repos {
		if err := runCommand("composer", QuietMode, "config", "repositories."+repo.Name, "composer", "https://"+repo.Url); err != nil {
			return fmt.Errorf("failed to add flux ui repository: %w", err)
		}
	}

	prod := p.config.PHP.Prod
	if err := runCommand("composer", QuietMode, append([]string{"require", "--no-update"}, prod...)...); err != nil {
		return fmt.Errorf("failed to add composer dependencies to composer.json (prod): %w", err)
	}

	dev := p.config.PHP.Dev

	if err := runCommand("composer", QuietMode, append([]string{"require", "--dev", "--no-update"}, dev...)...); err != nil {
		return fmt.Errorf("failed to add composer dependencies to composer.json (dev): %w", err)
	}

	if err := runCommand("composer", NormalMode, "update", "--no-scripts", "--no-interaction"); err != nil {
		return fmt.Errorf("failed to update composer dependencies: %w", err)
	}

	if err := runCommand("composer", NormalMode, "install"); err != nil {
		return fmt.Errorf("failed to install composer dependencies: %w", err)
	}

	artisanCmds := p.config.Artisan
	for _, cmd := range artisanCmds {
		parts := strings.Split(strings.TrimSpace(cmd), " ")
		if err := runCommand("php", InteractiveMode, append([]string{"artisan"}, parts...)...); err != nil {
			return fmt.Errorf("failed to run %s: %v", cmd, err)
		}
	}

	return nil
}
