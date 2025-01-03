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

//go:embed config.json
var configDotJson embed.FS

type AuthConfig struct {
	HTTPBasic map[string]HTTPBasicCredentials `json:"http-basic"`
}

type HTTPBasicCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PackageOptions struct {
	Prod   []string `json:"prod"`
	Dev    []string `json:"dev"`
	Remove []string `json:"remove"`
}

type Options struct {
	PHP     PackageOptions `json:"php"`
	JS      PackageOptions `json:"js"`
	Files   []string       `json:"files.remove"`
	Artisan []string       `json:"artisan"`
}

type config struct {
	projectName string
	projectDir  string
	opts        Options
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

	pkgConfig, err := getConfigDotJson()

	if err != nil {
		log.Fatal(err)
	}

	cfg := config{
		projectName: projectName,
		projectDir:  projectDir,
		opts:        pkgConfig,
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

	if err := runCommand("laravel", "new", cfg.projectName); err != nil {
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
		handleAuthJSON,
		handleCleanUp,
		handleGitignore,
		handlePHPDeps,
		handleJSDeps,
		handlePublishFiles,
	}
}

func getConfigDotJson() (Options, error) {
	data, err := configDotJson.ReadFile("config.json")

	if err != nil {
		return Options{}, fmt.Errorf("error reading embedded file: %w", err)
	}

	var cfg Options

	if err := json.Unmarshal(data, &cfg); err != nil {
		return Options{}, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return cfg, nil
}
