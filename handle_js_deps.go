package main

import (
	"fmt"
	"strings"
)

func handleJSDeps(p *project) error {
	fmt.Println("Installing node dependencies...")
	prod := p.config.JS.Prod
	if err := runQuietly("npm", append([]string{"install", "--package-lock-only"}, prod...)...); err != nil {
		return fmt.Errorf("failed to add node dependencies to package.json (prod): %w", err)
	}

	dev := p.config.JS.Dev
	if err := runQuietly("npm", append([]string{"install", "-D", "--package-lock-only"}, dev...)...); err != nil {
		return fmt.Errorf("failed to add node dependencies to package.json (dev): %w", err)
	}

	if err := run("npm", "install"); err != nil {
		return fmt.Errorf("failed to install node dependencies: %w", err)
	}

	npxCmds := p.config.NPX
	for _, cmd := range npxCmds {
		parts := strings.Split(strings.TrimSpace(cmd), " ")
		if err := runInteractive("npx", parts...); err != nil {
			return fmt.Errorf("failed to run %s: %v", cmd, err)
		}
	}

	return nil
}
