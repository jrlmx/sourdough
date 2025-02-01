package main

import (
	"fmt"
	"log"
	"strings"
)

type StarterConfigJson struct {
	Template string   `json:"template"`
	Flags    []string `json:"flags"`
	Composer struct {
		Repositories []struct {
			Name string `json:"name"`
			Url  string `json:"url"`
			Auth bool   `json:"auth"`
		} `json:"repositories"`
		Production  []string `json:"production"`
		Development []string `json:"development"`
		Remove      []string `json:"remove"`
	} `json:"composer"`
	Npm struct {
		Production  []string `json:"production"`
		Development []string `json:"development"`
		Remove      []string `json:"remove"`
	} `json:"npm"`
	Remove          []string `json:"remove_files"`
	Commands        []string `json:"commands"`
	ArtisanCommands []string `json:"artisan_commands"`
	NpxCommands     []string `json:"npx_commands"`
}

type config struct {
	args    cliOptions
	wd      string
	project *project
	starter *starter
	cm      cleanUpManager
}

type project struct {
	name string
	dir  string
}

type repo struct {
	name string
	url  string
	auth bool
}

type dependencyList struct {
	production  []string
	development []string
	remove      []string
	repos       *[]repo
}

type starter struct {
	template string
	flags    []string
	name     string
	dir      string
	composer dependencyList
	npm      dependencyList
	remove   []string
	commands map[string][]command
}

func (sj *StarterConfigJson) parse(name, dir string) (starter, error) {
	// Parse commands
	commands, err := sj.parseCommands()
	if err != nil {
		return starter{}, err
	}

	// Parse composer repositories
	var repos []repo
	for _, r := range sj.Composer.Repositories {
		if r.Name == "" {
			return starter{}, fmt.Errorf("composer repo name cannot be empty: %s", r.Name)
		}
		if r.Url == "" {
			return starter{}, fmt.Errorf("composer repo url cannot be empty: %s", r.Name)
		}
		repos = append(repos, repo{
			name: r.Name,
			url:  r.Url,
			auth: r.Auth,
		})
	}

	fmt.Print(sj.Flags)

	return starter{
		name:     name,
		dir:      dir,
		template: sj.Template,
		remove:   sj.Remove,
		composer: dependencyList{
			production:  sj.Composer.Production,
			development: sj.Composer.Development,
			remove:      sj.Composer.Remove,
			repos:       &repos,
		},
		npm: dependencyList{
			production:  sj.Npm.Production,
			development: sj.Npm.Development,
			remove:      sj.Npm.Remove,
		},
		flags:    sj.Flags,
		commands: commands,
	}, nil
}

func (sj *StarterConfigJson) parseCommands() (map[string][]command, error) {
	artisanCommands := []command{}
	for _, cstr := range sj.ArtisanCommands {
		parts := strings.Split(cstr, " ")
		cmd := newCommand("php", append([]string{"artisan"}, parts...), InteractiveMode)
		if err := cmd.validate(); err != nil {
			return nil, err
		}
		artisanCommands = append(artisanCommands, *cmd)
	}
	npxCommands := []command{}
	for _, cstr := range sj.NpxCommands {
		parts := strings.Split(cstr, " ")
		cmd := newCommand("npx", parts, InteractiveMode)
		if err := cmd.validate(); err != nil {
			return nil, err
		}
		npxCommands = append(npxCommands, *cmd)
	}
	commands := map[string][]command{
		"default":          {},
		"composer_install": artisanCommands,
		"npm_install":      npxCommands,
	}
	hooks := getHooks()
	for _, hook := range hooks {
		if _, ok := commands[hook]; !ok {
			commands[hook] = []command{}
		}
	}
	for _, cstr := range sj.Commands {
		hook := "default"
		mode := NormalMode
		if strings.HasPrefix(cstr, "@") {
			parts := strings.SplitN(cstr, ":", 2)
			hook = parts[0]
			cstr = parts[1]
		}
		if _, ok := commands[hook]; !ok {
			log.Fatal("invalid hook:", hook)
		}
		if strings.HasPrefix(cstr, "quiet:") {
			mode = QuietMode
			cstr = strings.TrimPrefix(cstr, "quiet:")
		} else if strings.HasPrefix(cstr, "interactive:") {
			mode = InteractiveMode
			cstr = strings.TrimPrefix(cstr, "interactive:")
		}
		parts := strings.Split(cstr, " ")
		cmd := newCommand(parts[0], parts[1:], mode)
		if err := cmd.validate(); err != nil {
			return nil, err
		}
		commands[hook] = append(commands[hook], *cmd)
	}
	return commands, nil
}
