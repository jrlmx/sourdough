package main

import (
	"fmt"
	"strings"
)

func handlePHPDeps(p *project) error {
	fmt.Println("Installing composer dependencies...")
	for _, repo := range p.config.Repos {
		if err := runQuietly("composer", "config", "repositories."+repo.Name, "composer", "https://"+repo.Url); err != nil {
			return fmt.Errorf("failed to add flux ui repository: %w", err)
		}
	}

	prod := p.config.PHP.Prod
	if err := runQuietly("composer", append([]string{"require", "--no-install"}, prod...)...); err != nil {
		return fmt.Errorf("failed to add composer dependencies to composer.json (prod): %w", err)
	}

	dev := p.config.PHP.Dev
	if err := runQuietly("composer", append([]string{"require", "--dev", "--no-install"}, dev...)...); err != nil {
		return fmt.Errorf("failed to add composer dependencies to composer.json (dev): %w", err)
	}

	if err := run("composer", "install"); err != nil {
		return fmt.Errorf("failed to install composer dependencies: %w", err)
	}

	artisanCmds := p.config.Artisan
	for _, cmd := range artisanCmds {
		parts := strings.Split(strings.TrimSpace(cmd), " ")
		if err := runInteractive("php", append([]string{"artisan"}, parts...)...); err != nil {
			return fmt.Errorf("failed to run %s: %v", cmd, err)
		}
	}

	return nil
}
