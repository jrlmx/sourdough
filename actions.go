package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func createLaravelProject(cfg *config) error {
	if err := exec.Command("which", "laravel").Run(); err != nil {
		return fmt.Errorf(("laravel installer is not installed"))
	}

	if _, err := os.Stat(cfg.projectName); err == nil {
		return fmt.Errorf("directory already exists: %s", cfg.projectDir)
	}

	fmt.Printf("Creating new Laravel project: %s\n", cfg.projectName)

	cmd := exec.Command("laravel", "new", cfg.projectName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create Laravel project: %w", err)
	}

	if err := os.Chdir(cfg.projectDir); err != nil {
		return fmt.Errorf("failed to change to project directory: %w", err)
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

	path := filepath.Join(cfg.projectDir, "auth.json")

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
	path := filepath.Join(cfg.projectDir, ".gitignore")

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

func installComposerDeps(cfg *config) error {
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
		return fmt.Errorf("failed to install composer dependencies: %w", err)
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

func installNodeDeps(cfg *config) error {
	if err := exec.Command("which", "npm").Run(); err != nil {
		return fmt.Errorf("npm is not installed")
	}

	deps := []string{
		"prettier",
		"prettier-plugin-blade",
		"@tailwindcss/typography",
		"@tailwindcss/forms",
	}

	cmd := exec.Command("npm", append([]string{"install", "-D"}, deps...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install node dependencies: %w", err)
	}

	return nil
}

func copyStubFiles(cfg *config) error {
	err := os.MkdirAll(cfg.projectDir, os.ModePerm)

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
		destPath := filepath.Join(cfg.projectDir, relPath)

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
