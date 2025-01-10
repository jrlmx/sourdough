package main

import (
	"fmt"
)

func handleJSDeps(p *project) error {
	prod := p.config.JS.Prod
	if err := run("npm", append([]string{"install", "--package-lock-only"}, prod...)...); err != nil {
		return fmt.Errorf("failed to add node dependencies to package.json (prod): %w", err)
	}

	dev := p.config.JS.Dev
	if err := run("npm", append([]string{"install", "-D", "--package-lock-only"}, dev...)...); err != nil {
		return fmt.Errorf("failed to add node dependencies to package.json (dev): %w", err)
	}

	if err := run("npm", "install"); err != nil {
		return fmt.Errorf("failed to install node dependencies: %w", err)
	}

	return nil
}
