package main

import (
	"fmt"
	"os"
	"slices"
)

func createProjectAction(cfg *config) error {
	fmt.Println("Creating project...")
	// Set the mode based on the interactive flag in the starter config
	mode := InteractiveMode
	if slices.ContainsFunc(cfg.starter.flags, func(s string) bool {
		return s == "-n" || s == "--no-interaction" || s == "--quiet" || s == "-q"
	}) {
		mode = QuietMode
	}
	// Add a cleanup task to remove the project directory
	cfg.cm.addTask(func() error {
		return os.RemoveAll(cfg.project.dir)
	}, false)
	// Start the installation process
	if cfg.starter.template != "" {
		// Clone the starter template if set
		fmt.Println("Cloning github repository:", cfg.starter.template)
		if err := runCommand("git", append([]string{"clone", "--depth=1", cfg.starter.template, cfg.project.name}, cfg.starter.flags...), mode); err != nil {
			return fmt.Errorf("error cloning repository: %w", err)
		}
	} else {
		// Use the Laravel installer to create a new project
		fmt.Println("Using Laravel Installer...")
		if err := runCommand("laravel", append([]string{"new", cfg.project.name}, cfg.starter.flags...), mode); err != nil {
			return fmt.Errorf("error creating project: %w", err)
		}
	}
	// Change directory to the project directory
	fmt.Println("Navigating to project directory...")
	if err := os.Chdir(cfg.project.dir); err != nil {
		return fmt.Errorf("error changing directory: %w", err)
	}
	return nil
}
