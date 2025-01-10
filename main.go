package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

//go:embed all:starters/*
var starters embed.FS

func main() {
	var name string

	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s [project_name]\n", os.Args[0])
	}

	name = os.Args[1]

	if _, err := os.Stat(name); err == nil {
		log.Fatalf("project directory already exists: %s", name)
	}

	cleaned := cleanDirName(name)

	wdir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	pdir := filepath.Join(wdir, cleaned)

	p := newProject(cleaned, pdir)

	if err := checkSystemDeps(); err != nil {
		log.Fatal(err)
	}

	for _, action := range actions() {
		if err := action(&p); err != nil {
			cleanupOnFailure(&p)
			log.Fatal(err)
		}
	}

	fmt.Println("Laravel project scaffolding complete!")
}

func actions() []func(p *project) error {
	return []func(p *project) error{
		handleCreateProject,
		handleStarterSelection,
		handleAuthJSON,
		handleCleanUp,
		handleGitignore,
		handlePHPDeps,
		handleJSDeps,
		handlePublishFiles,
	}
}

func cleanupOnFailure(p *project) error {
	fmt.Println("Cleaning up...")
	if _, err := os.Stat(p.pdir); err == nil {
		if err := os.RemoveAll(p.pdir); err != nil {
			return fmt.Errorf("failed to cleanup directory: %w", err)
		}
	}
	return nil
}

func checkSystemDeps() error {
	deps := []string{
		"php",
		"laravel",
		"composer",
		"npm",
	}

	for _, dep := range deps {
		path, err := isInstalled(dep)
		if err != nil {
			return err
		}
		fmt.Printf("Using %s: %s\n", dep, path)
	}

	return nil
}
