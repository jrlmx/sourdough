package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

//go:embed stubs/*
var stubs embed.FS

//go:embed config.json
var configDotJson embed.FS

type Options struct {
	PHP     PackageOptions `json:"php"`
	JS      PackageOptions `json:"js"`
	Files   []string       `json:"files.remove"`
	Artisan []string       `json:"artisan"`
}

type PackageOptions struct {
	Prod   []string `json:"prod"`
	Dev    []string `json:"dev"`
	Remove []string `json:"remove"`
}

type project struct {
	name string
	dir  string
	opts Options
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: program <project-name>")
	}

	wd, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	name := os.Args[1]
	dir := filepath.Join(wd, name)

	if _, err := os.Stat(name); err == nil {
		log.Fatalf("directory already exists: %s", dir)
	}

	opts, err := getConfigDotJson()

	if err != nil {
		log.Fatal(err)
	}

	p := project{
		name,
		dir,
		opts,
	}

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

func getConfigDotJson() (Options, error) {
	data, err := configDotJson.ReadFile("config.json")

	if err != nil {
		return Options{}, fmt.Errorf("error reading embedded file: %w", err)
	}

	var p Options

	if err := json.Unmarshal(data, &p); err != nil {
		return Options{}, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return p, nil
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
