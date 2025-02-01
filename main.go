package main

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
)

//go:embed config.json
var fsys embed.FS

type cliOptions struct {
	projectName string
	starterName string
	printHooks  bool
}

func getActions() []Action {
	return []Action{
		{name: "create_project", hookable: false, callback: createProjectAction},
		{name: "create_auth_json", hookable: false, callback: createAuthJsonAction},
		{name: "remove_files", hookable: false, callback: removeFilesAction},
		{name: "composer_install", hookable: true, callback: composerInstallAction},
		{name: "npm_install", hookable: true, callback: npmInstallAction},
		{name: "publish_files", hookable: true, callback: publishFilesAction},
		{name: "run_commands", hookable: false, callback: runCommandsAction},
	}
}

func getHooks() []string {
	var hooks []string
	for _, a := range getActions() {
		if a.hookable {
			hooks = append(hooks, a.name)
		}
	}
	return hooks
}

func main() {
	cliOpts := getCliOptions()

	if cliOpts.printHooks {
		fmt.Println("Available command hooks:")
		for _, hook := range getHooks() {
			fmt.Println(hook)
		}
		fmt.Println("Note: hooks execute after the specified action")
		os.Exit(0)
	}

	sc := make(chan os.Signal, 1)
	ec := make(chan error, 1)
	signal.Notify(sc, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	cfg := config{
		args:    cliOpts,
		project: &project{},
		starter: &starter{},
		wd:      wd,
	}

	go func() {
		if err := run(&cfg); err != nil {
			ec <- err
		}
		cancel()
	}()

	select {
	case <-sc:
		fmt.Printf("Interrupted")
		cfg.cm.cleanUp(false)
	case err := <-ec:
		fmt.Println(err)
		cfg.cm.cleanUp(false)
	case <-ctx.Done():
		fmt.Println("Done")
		cfg.cm.cleanUp(true)
	}
	cancel()
}

func run(cfg *config) error {
	// Get the sourdough config
	sdConfig, err := getSourdoughConfig()
	if err != nil {
		return err
	}
	// Name the project
	err = cfg.createProjectConfig(cfg.args.projectName)
	if err != nil {
		return err
	}
	// Select a starter
	err = cfg.createStarterConfig(cfg.args.starterName, sdConfig.Starters)
	if err != nil {
		return err
	}
	// Ensure config project and starter are set
	if cfg.project == nil || cfg.starter == nil {
		return errors.New("project or starter is nil")
	}
	// Run the actions
	for _, a := range getActions() {
		if err := a.callback(cfg); err != nil {
			return err
		}
		if _, ok := cfg.starter.commands[a.name]; ok {
			for _, cmd := range cfg.starter.commands[a.name] {
				if err := cmd.run(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func getCliOptions() cliOptions {
	var projectName string
	var starterName string
	var printHooks bool

	flag.StringVar(&starterName, "starter", "", "Name of the starter")
	flag.BoolVar(&printHooks, "hooks", false, "Print available hooks")
	flag.Parse()

	if len(flag.Args()) > 0 {
		projectName = flag.Args()[0]
	}

	return cliOptions{
		projectName: projectName,
		starterName: starterName,
	}
}

func getSourdoughConfig() (SourdoughConfig, error) {
	file, err := fsys.ReadFile("config.json")
	if err != nil {
		return SourdoughConfig{}, fmt.Errorf("error reading config.json: %w", err)
	}
	var cfg SourdoughConfig
	if err := json.Unmarshal(file, &cfg); err != nil {
		return SourdoughConfig{}, fmt.Errorf("error parsing config.json: %w", err)
	}
	if len(cfg.Starters) < 1 {
		return SourdoughConfig{}, errors.New("no starters found in config.json")
	}
	return cfg, err
}

func (cfg *config) createProjectConfig(name string) error {
	cleaned := strings.TrimSpace(name)
	isValid := func(s string) error {
		if s == "" {
			return errors.New("project name cannot be empty")
		}
		if len(s) < 1 || len(name) >= 60 {
			return errors.New("project name must be between 1 and 60 characters")
		}
		if strings.ContainsAny(s, "&;|<>(){}$`\\/.") {
			return errors.New("project name contains invalid characters")
		}
		return nil
	}
	if cleaned == "" || isValid(cleaned) != nil {
		if err := huh.NewInput().
			Title("Name your project").
			Value(&cleaned).
			Validate(isValid).
			Run(); err != nil {
			return fmt.Errorf("error naming project '%s': %w", cleaned, err)
		}
	}
	dir := filepath.Join(cfg.wd, cleaned)
	// Ensure the project directory does not already exist
	if _, err := os.Stat(dir); err == nil {
		return errors.New("project directory already exists")
	}
	cfg.project = &project{
		name: cleaned,
		dir:  dir,
	}
	return nil
}

func (cfg *config) createStarterConfig(name string, options map[string]string) error {
	cleaned := strings.TrimSpace(name)
	isValid := func(s string) error {
		if s == "" {
			return errors.New("starter name cannot be empty")
		}
		if _, ok := options[s]; !ok {
			return errors.New("starter name not found in config.starters")
		}
		return nil
	}
	if cleaned == "" || isValid(cleaned) != nil {
		selectOptions := []huh.Option[string]{}
		for key := range options {
			selectOptions = append(selectOptions, huh.NewOption(key, key))
		}
		if err := huh.NewSelect[string]().
			Options(selectOptions...).
			Title("Select a starter").
			Value(&cleaned).
			Validate(isValid).
			Run(); err != nil {
			return fmt.Errorf("error selecting starter '%s': %w", cleaned, err)
		}
	}
	// Create a temp directory for the starter
	dir, err := os.MkdirTemp(cfg.wd, "sd-tmp-")
	if err != nil {
		return fmt.Errorf("error creating temp directory for starter '%s': %w", cleaned, err)
	}
	cfg.cm.addTask(func() error {
		return os.RemoveAll(dir)
	}, true)
	// Clone the git repo into the temp directory
	if err := runCommand("git", []string{"clone", options[cleaned], dir}, QuietMode); err != nil {
		return fmt.Errorf("error cloning starter '%s': %w", cleaned, err)
	}
	// Parse the starter config
	var configJson StarterConfigJson
	file, err := os.ReadFile(filepath.Join(dir, "config.json"))
	if err != nil {
		return fmt.Errorf("error reading starter config '%s': %w", cleaned, err)
	}
	if err := json.Unmarshal(file, &configJson); err != nil {
		return fmt.Errorf("error parsing starter config '%s': %w", cleaned, err)
	}
	// Create the starter
	starter, err := configJson.parse(cleaned, dir)
	if err != nil {
		return fmt.Errorf("error creating starter '%s': %w", cleaned, err)
	}
	cfg.starter = &starter
	return nil
}
