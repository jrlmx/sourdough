package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed stubs/*
var stubs embed.FS

//go:embed deps.json
var depsDotJson embed.FS

type AuthConfig struct {
	HTTPBasic map[string]HTTPBasicCredentials `json:"http-basic"`
}

type HTTPBasicCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PkgConfig struct {
	Composer []string `json:"composer"`
	NPM      []string `json:"npm"`
}

type config struct {
	projectName string
	projectDir  string
	flux        *struct {
		username string
		password string
	}
	deps PkgConfig
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: program <project-name>")
	}

	workingDir, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	projectName := os.Args[1]
	projectDir := filepath.Join(workingDir, projectName)

	pkgConfig, err := getPkgConfig()

	if err != nil {
		log.Fatal(err)
	}

	cfg := config{
		projectName: projectName,
		projectDir:  projectDir,
		deps:        pkgConfig,
	}

	err = createApp(&cfg)

	if err != nil {
		log.Fatal(err)
	}

	if err := applyStarterKit(&cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Laravel project scaffolding complete!")
}

func applyStarterKit(cfg *config) error {
	for _, action := range getActions() {
		if err := action(cfg); err != nil {
			cleanupOnFailure(cfg)
			return fmt.Errorf("failed to apply starter kit: %w", err)
		}
	}
	return nil
}

func createApp(cfg *config) error {
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

func cleanupOnFailure(cfg *config) error {
	if _, err := os.Stat(cfg.projectName); err == nil {
		if err := os.RemoveAll(cfg.projectName); err != nil {
			return fmt.Errorf("failed to cleanup directory: %w", err)
		}
	}
	return nil
}

func getActions() []func(cfg *config) error {
	return []func(cfg *config) error{
		handleFluxPrompt,
		handleAuthJSON,
		handleGitignore,
		handleComposerDeps,
		handleNodeDeps,
		handlePublishFiles,
	}
}

func getPkgConfig() (PkgConfig, error) {
	data, err := depsDotJson.ReadFile("deps.json")

	if err != nil {
		return PkgConfig{}, fmt.Errorf("error reading embedded file: %w", err)
	}

	var cfg PkgConfig

	if err := json.Unmarshal(data, &cfg); err != nil {
		return PkgConfig{}, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return cfg, nil
}
