package main

import (
	"fmt"
	"os"
)

func handleCreateProject(p *project) error {
	fmt.Printf("Creating new Laravel project: %s\n", p.name)

	if err := run("laravel", "new", p.name); err != nil {
		return fmt.Errorf("failed to create Laravel project: %w", err)
	}

	if err := os.Chdir("./" + p.name); err != nil {
		return fmt.Errorf("failed to change to project directory: %w", err)
	}

	return nil
}
