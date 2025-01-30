package main

import (
	"fmt"
	"os"
	"slices"
)

func createProjectAction(c *config) error {
	fmt.Println("Creating project...")
	// Set the mode based on the interactive flag in the starter config
	mode := InteractiveMode
	print("c.starter.flags: ", c.starter.flags)
	if slices.ContainsFunc(c.starter.flags, func(s string) bool {
		return s == "-n" || s == "--no-interaction" || s == "--quiet" || s == "-q"
	}) {
		mode = QuietMode
	}
	// Add a cleanup task to remove the project directory
	c.cm.addTask(func() error {
		return os.RemoveAll(c.project.dir)
	}, false)
	// Start the installation process
	if c.starter.template != "" {
		// Clone the starter template if set
		fmt.Println("Cloning github repository:", c.starter.template)
		if err := runCommand("git", append([]string{"clone", "--depth=1", c.starter.template, c.project.name}, c.starter.flags...), mode); err != nil {
			return fmt.Errorf("error cloning repository: %w", err)
		}
	} else {
		// Use the Laravel installer to create a new project
		fmt.Println("Using Laravel Installer...")
		if err := runCommand("laravel", append([]string{"new", c.project.name}, c.starter.flags...), mode); err != nil {
			return fmt.Errorf("error creating project: %w", err)
		}
	}
	// Change directory to the project directory
	fmt.Println("Navigating to project directory...")
	if err := os.Chdir(c.project.dir); err != nil {
		return fmt.Errorf("error changing directory: %w", err)
	}
	return nil
}
