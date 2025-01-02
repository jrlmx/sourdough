package main

import (
	"fmt"
	"os"
	"os/exec"
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
		err := os.Remove(projectDir + "/" + file)

		if err != nil {
			return err
		}
	}

	return nil
}
