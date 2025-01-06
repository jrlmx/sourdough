package main

import (
	"fmt"
)

func handlePHPDeps(p *project) error {
	fmt.Println("Installing composer dependencies...")

	if err := run("composer", "config", "repositories.flux-pro", "composer", "https://composer.fluxui.dev"); err != nil {
		return fmt.Errorf("failed to add flux ui repository: %w", err)
	}

	prod := p.opts.PHP.Prod

	if err := run("composer", append([]string{"require"}, prod...)...); err != nil {
		return fmt.Errorf("failed to install composer dependencies: %w", err)
	}

	dev := p.opts.PHP.Dev

	if err := run("composer", append([]string{"require", "--dev"}, dev...)...); err != nil {
		return fmt.Errorf("failed to install composer dependencies: %w", err)
	}

	artisanCmds := p.opts.Artisan

	for _, cmd := range artisanCmds {
		if err := run("php", append([]string{"artisan"}, cmd)...); err != nil {
			return fmt.Errorf("failed to run %s: %v", cmd, err)
		}
	}

	return nil
}
