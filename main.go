package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//go:embed all:starters/*
var starters embed.FS

func main() {
	kitFlag := flag.String("kit", "", "Specify the starter kit to use")
	configFlag := flag.String("config", "", "Output the embeded config file for the specified kit and exit")
	treeFlag := flag.String("tree", "", "Output the embeded file tree for the specified kit and exit")

	flag.Parse()

	if *configFlag != "" {
		if err := printConfig(*configFlag); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	if *treeFlag != "" {
		if err := printFileTree(*treeFlag); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	var name string
	args := flag.Args()
	if len(args) < 1 {
		log.Fatalf("Usage: %s [project_name]\n", os.Args[0])
	}
	name = args[1]

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
	if kitFlag != nil {
		p.kit = kitFlag
	}

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
	os.Exit(0)
}

func actions() []func(p *project) error {
	return []func(p *project) error{
		handleStarterSelection,
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

func printConfig(starter string) error {
	data, err := starters.ReadFile("starters/" + starter + "/config.json")
	if err != nil {
		return err
	}

	var config Config

	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	prettyPrinted, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	fmt.Println(string(prettyPrinted))
	return nil
}

func printFileTree(starter string) error {
	basePath := filepath.Join("starters", starter)

	err := fs.WalkDir(starters, basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relativePath := strings.TrimPrefix(path, basePath)
		depth := strings.Count(relativePath, "/")
		indent := strings.Repeat("    ", depth)

		if path == basePath {
			fmt.Printf("%s/\n", filepath.Base(basePath))
			return nil
		}

		if d.IsDir() {
			fmt.Printf("%s├── %s/\n", indent, filepath.Base(path))
		} else {
			fmt.Printf("%s├── %s\n", indent, d.Name())
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
