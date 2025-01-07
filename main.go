package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

//go:embed kits/*
var kits embed.FS

type Config struct {
	PHP     PackageManifest `json:"php"`
	JS      PackageManifest `json:"js"`
	Files   []string        `json:"remove_files"`
	Repos   []Repo          `json:"repos"`
	Artisan []string        `json:"artisan"`
}

type PackageManifest struct {
	Remove []string `json:"remove"`
	Prod   []string `json:"prod"`
	Dev    []string `json:"dev"`
}

type Repo struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Auth bool   `json:"auth"`
}

type project struct {
	name   string
	kit    *string
	config *Config
}

func newProject(name string) project {
	return project{
		name,
		nil,
		nil,
	}
}

func (p *project) loadConfig() error {
	data, err := kits.ReadFile("kits/" + *p.kit + "/config.json")
	if err != nil {
		return err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("error parsing config.json: %w", err)
	}

	p.config = &cfg
	return nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: program <project-name>")
	}

	name := os.Args[1]

	if _, err := os.Stat(name); err == nil {
		log.Fatalf("project directory already exists: %s", name)
	}

	cleaned := cleanString(name)

	p := newProject(cleaned)

	if err := checkDependencies(); err != nil {
		log.Fatal(err)
	}

	for _, action := range actions() {
		if err := action(&p); err != nil {
			cleanupOnFailure(&p)
			log.Fatal("failed to apply starter kit\n", err)
		}
	}

	fmt.Println("Laravel project scaffolding complete!")
}

func actions() []func(p *project) error {
	return []func(p *project) error{
		handleCreateProject,
		handleKitSelection,
		handleAuthJSON,
		handleCleanUp,
		handleGitignore,
		handlePHPDeps,
		handleJSDeps,
		handlePublishFiles,
	}
}

func cleanupOnFailure(p *project) error {
	if _, err := os.Stat(p.name); err == nil {
		if err := os.RemoveAll(p.name); err != nil {
			return fmt.Errorf("failed to cleanup directory: %w", err)
		}
	}
	return nil
}

func checkDependencies() error {
	deps := []string{
		"php",
		"laravel",
		"composer",
		"npm",
	}

	for _, dep := range deps {
		if err := check(dep); err != nil {
			return err
		}
	}

	return nil
}
