package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func handleCleanUp(p *project) error {
	err := cleanUpFiles(p.config.Files)
	if err != nil {
		return err
	}

	err = cleanComposerPackages(p.config.PHP.Remove)
	if err != nil {
		return err
	}

	err = cleanJSPackages(p.config.JS.Remove)
	if err != nil {
		return err
	}

	return nil
}

func cleanJSPackages(packages []string) error {
	if len(packages) == 0 {
		fmt.Println("No JS packages to clean up.")
		return nil
	}

	fmt.Println("Cleaning up unwanted node packages...")
	if err := run("npm", append([]string{"remove"}, packages...)...); err != nil {
		return fmt.Errorf("failed to remove node packages: %w", err)
	}

	return nil
}

func cleanComposerPackages(packages []string) error {
	if len(packages) == 0 {
		fmt.Println("No PHP packages to clean up.")
		return nil
	}

	fmt.Println("Cleaning up unwanted composer packages...")
	if err := run("composer", append([]string{"remove"}, packages...)...); err != nil {
		return fmt.Errorf("failed to remove composer packages: %w", err)
	}

	return nil
}

func cleanUpFiles(files []string) error {
	if len(files) == 0 {
		fmt.Println("No files to clean up.")
		return nil
	}

	fmt.Println("Cleaning up unwanted files...")

	for _, file := range files {
		cleanPath := filepath.Join(".", file)

		err := os.RemoveAll(cleanPath)
		if err != nil {
			return err
		}
	}

	return nil
}
