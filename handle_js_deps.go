package main

import (
	"fmt"
	"strings"
)

func handleJSDeps(p *project) error {
	fmt.Println("Installing node dependencies...")
	prod := p.config.JS.Prod
	if err := runCommand("npm", QuietMode, append([]string{"install", "--no-package-lock"}, prod...)...); err != nil {
		return fmt.Errorf("failed to add node dependencies to package.json (prod): %w", err)
	}

	dev := p.config.JS.Dev
	if err := runCommand("npm", QuietMode, append([]string{"install", "-D", "--no-package-lock"}, dev...)...); err != nil {
		return fmt.Errorf("failed to add node dependencies to package.json (dev): %w", err)
	}

	if err := runCommand("npm", NormalMode, "install"); err != nil {
		return fmt.Errorf("failed to install node dependencies: %w", err)
	}

	npxCmds := p.config.NPX
	for _, cmd := range npxCmds {
		parts := strings.Split(strings.TrimSpace(cmd), " ")
		if err := runCommand("npx", InteractiveMode, parts...); err != nil {
			return fmt.Errorf("failed to run %s: %v", cmd, err)
		}
	}

	return nil
}
