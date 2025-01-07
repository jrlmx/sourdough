package main

import (
	"embed"
	"fmt"
	"log"
	"os"
)

//go:embed starters/*
var starters embed.FS

func main() {
	var name string

	if len(os.Args) != 2 {
		log.Fatal("usage: program <project-name>")
	}

	name = os.Args[1]

	if _, err := os.Stat(name); err == nil {
		log.Fatalf("project directory already exists: %s", name)
	}

	cleaned := cleanString(name)

	p := newProject(cleaned)

	if err := checkDeps(); err != nil {
		log.Fatal(err)
	}

	for _, action := range actions() {
		if err := action(&p); err != nil {
			cleanupOnFailure(&p)
			log.Fatal("failed to apply starter starter\n", err)
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

func checkDeps() error {
	deps := []string{
		"php",
		"laravel",
		"composer",
		"npm",
	}

	for _, dep := range deps {
		path, err := check(dep)
		if err != nil {
			return err
		}
		fmt.Printf("Using %s: %s", dep, path)
	}

	return nil
}
