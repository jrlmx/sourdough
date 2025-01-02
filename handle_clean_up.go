package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func handleCleanUp(cfg *config) error {
	err := cleanUpFiles(cfg.projectDir, cfg.options.Cleanup.Files)

	if err != nil {
		return err
	}

	err = cleanComposerPackages(cfg.options.Cleanup.Packages.PHP)

	if err != nil {
		return err
	}

	err = cleanJSPackages(cfg.options.Cleanup.Packages.JS)

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

	cmd := exec.Command("npm", append([]string{"remove"}, packages...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
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

	cmd := exec.Command("composer", append([]string{"remove"}, packages...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove composer packages: %w", err)
	}

	return nil
}

func cleanUpFiles(projectDir string, files []string) error {
	if len(files) == 0 {
		fmt.Println("No files to clean up.")
		return nil
	}

	fmt.Println("Cleaning up unwanted files...")

	for _, file := range files {
		cleanPath := filepath.Clean(projectDir + "/" + file)

		if strings.Contains(cleanPath, "..") {
			return fmt.Errorf("parent directory traversal is not allowed: %s", cleanPath)
		}

		if !strings.Contains(cleanPath, projectDir) {
			return fmt.Errorf("file is not in project directory: %s", cleanPath)
		}

		err := os.Remove(cleanPath)

		if err != nil {
			return err
		}
	}

	return nil
}
