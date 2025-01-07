package main

import (
	"fmt"
)

func handleJSDeps(p *project) error {
	prod := p.config.JS.Prod

	if err := run("npm", append([]string{"install"}, prod...)...); err != nil {
		return fmt.Errorf("failed to install node dependencies: %w", err)
	}

	dev := p.config.JS.Dev

	if err := run("npm", append([]string{"install", "-D"}, dev...)...); err != nil {
		return fmt.Errorf("failed to install node dependencies: %w", err)
	}

	return nil
}
