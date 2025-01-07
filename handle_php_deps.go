package main

import (
	"fmt"
)

func handlePHPDeps(p *project) error {
	fmt.Println("Installing composer dependencies...")

	for _, repo := range p.config.Repos {
		if err := run("composer", "config", "repositories."+repo.Name, "composer", "https://"+repo.Url); err != nil {
			return fmt.Errorf("failed to add flux ui repository: %w", err)
		}
	}

	prod := p.config.PHP.Prod

	if err := run("composer", append([]string{"require"}, prod...)...); err != nil {
		return fmt.Errorf("failed to install composer dependencies: %w", err)
	}

	dev := p.config.PHP.Dev

	if err := run("composer", append([]string{"require", "--dev"}, dev...)...); err != nil {
		return fmt.Errorf("failed to install composer dependencies: %w", err)
	}

	artisanCmds := p.config.Artisan

	for _, cmd := range artisanCmds {
		if err := run("php", append([]string{"artisan"}, cmd)...); err != nil {
			return fmt.Errorf("failed to run %s: %v", cmd, err)
		}
	}

	return nil
}
