package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

//go:embed stubs/*
var stubs embed.FS

type AuthConfig struct {
	HTTPBasic map[string]HTTPBasicCredentials `json:"http-basic"`
}

type HTTPBasicCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type config struct {
	projectName string
	projectDir  string
	flux        *struct {
		username string
		password string
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: program <project-name>")
	}

	projectName := os.Args[1]
	workingDir, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	projectDir := filepath.Join(workingDir, projectName)

	cfg := config{
		projectName: projectName,
		projectDir:  projectDir,
	}

	for _, action := range getActions() {
		if err := action(&cfg); err != nil {
			cleanupOnFailure(&cfg)
			log.Fatal(err)
		}
	}

	fmt.Println("Laravel project scaffolding complete!")
}

func getActions() []func(cfg *config) error {
	return []func(cfg *config) error{
		createLaravelProject,
		fluxPrompt,
		createAuthJSON,
		updateGitignore,
		installComposerDeps,
		installNodeDeps,
		copyStubFiles,
	}
}

func cleanupOnFailure(cfg *config) error {
	if _, err := os.Stat(cfg.projectName); err == nil {
		if err := os.RemoveAll(cfg.projectName); err != nil {
			return fmt.Errorf("failed to cleanup directory: %w", err)
		}
	}
	return nil
}
