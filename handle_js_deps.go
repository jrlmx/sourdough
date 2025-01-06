package main

import (
	"fmt"
)

func handleJSDeps(p *project) error {
	prod := p.opts.JS.Prod

	if err := run("npm", append([]string{"install"}, prod...)...); err != nil {
		return fmt.Errorf("failed to install node dependencies: %w", err)
	}

	dev := p.opts.JS.Dev

	if err := run("npm", append([]string{"install", "-D"}, dev...)...); err != nil {
		return fmt.Errorf("failed to install node dependencies: %w", err)
	}

	return nil
}
