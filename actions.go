package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func isLaravelProject(cfg *config) error {
	if _, err := os.Stat(filepath.Join(cfg.dir, "composer.json")); errors.Is(err, os.ErrNotExist) {
		return errors.New("current directory is not a laravel project (no composer.json found)")
	}

	return nil
}

func fluxPrompt(cfg *config) error {
	r := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Flux UI Username (email): ")
	username, _ := r.ReadString('\n')

	fmt.Print("Enter Flux UI License Key: ")
	license, _ := r.ReadString('\n')

	cfg.flux = &struct {
		username string
		password string
	}{
		username: strings.TrimSpace(username),
		password: strings.TrimSpace(license),
	}

	return nil
}

func createAuthJSON(cfg *config) error {
	authConfig := AuthConfig{
		HTTPBasic: map[string]HTTPBasicCredentials{
			"composer.fluxui.dev": {
				Username: cfg.flux.username,
				Password: cfg.flux.password,
			},
		},
	}

	path := filepath.Join(cfg.dir, "auth.json")

	file, err := json.MarshalIndent(authConfig, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to create auth.json: %w", err)
	}

	if err := os.WriteFile(path, file, 0644); err != nil {
		return fmt.Errorf("failed to write auth.json: %w", err)
	}

	return nil
}

func updateGitignore(cfg *config) error {
	path := filepath.Join(cfg.dir, ".gitignore")

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read .gitignore: %w", err)
	}

	if strings.Contains(string(content), "auth.json") {
		fmt.Println("skipped: adding auth.json to .gitignore - already present")
		return nil
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open .gitignore: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString("\nauth.json\n"); err != nil {
		return fmt.Errorf("failed to update .gitignore: %w", err)
	}

	fmt.Println("Added auth.json to .gitignore")

	return nil
}

func installComposerDependencies(cfg *config) error {
	if err := exec.Command("which", "composer").Run(); err != nil {
		return fmt.Errorf("composer is not installed")
	}

	if err := exec.Command("composer", "config", "repositories.flux-pro", "composer", "https://composer.fluxui.dev").Run(); err != nil {
		return fmt.Errorf("failed to add flux ui repository: %w", err)
	}

	deps := []string{
		"livewire/livewire",
		"livewire/volt",
		"laravel/folio",
		"livewire/flux",
		"livewire/flux-pro",
	}

	cmd := exec.Command("composer", append([]string{"require"}, deps...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install Composer dependencies: %w", err)
	}

	artisanCmds := []string{
		"folio:install",
		"volt:install",
	}

	for _, artisanCmd := range artisanCmds {
		cmd := exec.Command("php", "artisan", artisanCmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Failed to run %s: %v", artisanCmd, err)
		}
	}

	return nil
}

func copyStubFiles(cfg *config) error {
	err := os.MkdirAll(cfg.dir, os.ModePerm)

	if err != nil {
		return fmt.Errorf("failed to create desination directory: %w", err)
	}

	err = fs.WalkDir(stubs, "stubs", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking directory: %w", err)
		}

		if d.IsDir() {
			return nil
		}

		content, err := stubs.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", path, err)
		}

		relPath, err := filepath.Rel("stubs", path)
		if err != nil {
			return fmt.Errorf("failed to determine relative path: %w", err)
		}
		destPath := filepath.Join(cfg.dir, relPath)

		err = os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", destPath, err)
		}

		err = os.WriteFile(destPath, content, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", destPath, err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	fmt.Println("Files copied successfully.")

	return nil
}
