package main

import (
	"fmt"
	"os"
)

func handleCreateProject(p *project) error {
	fmt.Printf("Creating new Laravel project: %s\n", p.name)

	if p.config.GitURL != "" {
		if err := runCommand("git", QuietMode, "clone", p.config.GitURL, p.name); err != nil {
			return fmt.Errorf("failed to clone starter template: %w", err)
		}
	} else {
		if err := runCommand("laravel", QuietMode, append([]string{"new", p.name}, p.config.Flags...)...); err != nil {
			return fmt.Errorf("failed to create Laravel project: %w", err)
		}
	}

	if err := os.Chdir("./" + p.name); err != nil {
		return fmt.Errorf("failed to change to project directory: %w", err)
	}

	return nil
}
