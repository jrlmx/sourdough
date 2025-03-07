package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
)

func getActions() []Action {
	return []Action{
		{name: "create_project", hookable: false, callback: createProjectAction},
		{name: "remove_files", hookable: false, callback: removeFilesAction},
		{name: "create_auth_json", hookable: false, callback: createAuthJsonAction},
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

	handleSubCommands(&cliOpts)

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

func handleSubCommands(cliOpts *cliOptions) {
	if cliOpts.printHooks {
		fmt.Println("Available command hooks:")
		for _, hook := range getHooks() {
			fmt.Println(hook)
		}
		fmt.Println("Note: hooks execute after the specified action")
		os.Exit(0)
	}

	if cliOpts.addStarter {
		if err := addNewStarter(); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Starter added successfully - run 'sourdough' again to use it\n")
		os.Exit(0)
	}

	if cliOpts.removeStarter != "" {
		if err := removeStarter(cliOpts.removeStarter); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Starter removed successfully\n")
		os.Exit(0)
	}

	if cliOpts.exportStarters {
		sdConfig, err := getSourdoughConfig()
		if err != nil {
			log.Fatal(err)
		}
		if err := exportStarters(sdConfig); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	if cliOpts.importStarters != "" {
		if err := importStarters(cliOpts.importStarters); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}
}

func getSourdoughConfig() (SourdoughConfig, error) {
	db, err := database()
	if err != nil {
		return SourdoughConfig{}, fmt.Errorf("error opening database: %w", err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT name, url FROM starters")
	if err != nil {
		return SourdoughConfig{}, fmt.Errorf("error querying database: %w", err)
	}
	starters := map[string]string{}
	for rows.Next() {
		var name, url string
		if err := rows.Scan(&name, &url); err != nil {
			return SourdoughConfig{}, fmt.Errorf("error scanning database row: %w", err)
		}
		starters[name] = url
	}
	if len(starters) < 1 {
		return SourdoughConfig{}, errors.New("no starters found in database")
	}
	return SourdoughConfig{Starters: starters}, nil
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

func addNewStarter() error {
	var name, url, description string
	if err := huh.NewForm(huh.NewGroup(
		huh.NewInput().
			Title("Name").
			Description("Enter the name of the starter").
			Value(&name).
			Validate(huh.ValidateNotEmpty()),
		huh.NewInput().
			Title("URL").
			Description("Enter the URL of the starter repository (include '.git')").
			Value(&url).
			Validate(func(s string) error {
				if err := huh.ValidateNotEmpty()(s); err != nil {
					return err
				}
				if !strings.HasPrefix(s, "https://") && !strings.HasPrefix(s, "git@") && !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "/") {
					return errors.New("URL must start with 'https://', 'git@', 'http://', or '/'")
				}
				if !strings.HasSuffix(s, ".git") {
					return errors.New("URL must end with '.git'")
				}
				return nil
			}),
		huh.NewText().
			Title("Description").
			Description("Optional: Enter a description of the starter").
			Value(&description),
	)).Run(); err != nil {
		return fmt.Errorf("error getting starter name and url: %w", err)
	}
	db, err := database()
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO starters (name, url, description) VALUES (?, ?, ?)", name, url, description)
	if err != nil {
		return fmt.Errorf("error inserting starter into database: %w", err)
	}
	return nil
}

func removeStarter(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("name cannot be empty")
	}
	db, err := database()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM starters WHERE name = ?", name)
	if err != nil {
		return err
	}
	return nil
}

func exportStarters(sdConfig SourdoughConfig) error {
	file, err := os.Create("sourdough.json")
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(sdConfig); err != nil {
		return err
	}
	fmt.Println("Exported starters to sourdough.json")
	return nil
}

func importStarters(path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return errors.New("path cannot be empty (i.e. ./sourdough.json)")
	}
	if !strings.HasSuffix(path, ".json") {
		return errors.New("path must point to a JSON file (i.e. sourdough.json)")
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	var sdConfig SourdoughConfig
	dec := json.NewDecoder(file)
	if err := dec.Decode(&sdConfig); err != nil {
		return err
	}
	db, err := database()
	if err != nil {
		return err
	}
	defer db.Close()
	for name, url := range sdConfig.Starters {
		_, err := db.Exec("INSERT INTO starters (name, url) VALUES (?, ?)", name, url)
		if err != nil {
			return err
		}
	}
	fmt.Println("Imported starters from sourdough.json")
	return nil
}
