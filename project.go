package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Config struct {
	GitURL   string          `json:"git_url"`
	PHP      PackageManifest `json:"php"`
	JS       PackageManifest `json:"js"`
	Files    []string        `json:"remove_files"`
	Repos    []Repo          `json:"repos"`
	Artisan  []string        `json:"artisan"`
	NPX      []string        `json:"npx"`
	Commands []string        `json:"commands"`
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
	name     string
	pdir     string
	kit      *string
	config   *Config
	commands map[string][]command
}

type command struct {
	name string
	args []string
	mode runMode
}

func newProject(name, pdir string, hooks []string) project {
	commands := map[string][]command{
		"default": {},
	}

	for _, hook := range hooks {
		commands[hook] = []command{}
	}

	return project{
		name,
		pdir,
		nil,
		nil,
		commands,
	}
}

func (p *project) init(kit string) error {
	p.kit = &kit

	if err := p.loadConfig(); err != nil {
		return err
	}

	if err := p.parseCommands(); err != nil {
		return err
	}

	return nil
}

func (p *project) loadConfig() error {
	data, err := starters.ReadFile("starters/" + *p.kit + "/config.json")
	if err != nil {
		return fmt.Errorf("failed to read config.json: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse config.json: %w", err)
	}

	p.config = &cfg
	return nil
}

func (p *project) parseCommands() error {
	for _, cmd := range p.config.Commands {
		hook := "default"
		mode := NormalMode

		if strings.HasPrefix(cmd, "@") {
			parts := strings.SplitN(cmd, ":", 2)
			hook = parts[0]
			cmd = parts[1]
		}

		if strings.HasPrefix(cmd, "quiet:") {
			mode = QuietMode
			cmd = strings.TrimPrefix(cmd, "quiet:")
		} else if strings.HasPrefix(cmd, "interact:") {
			mode = InteractiveMode
			cmd = strings.TrimPrefix(cmd, "interact:")
		}

		parts := strings.Split(cmd, " ")

		if _, ok := p.commands[hook]; !ok {
			return fmt.Errorf("invalid command hook %s in cmd: %s", hook, cmd)
		}

		p.commands[hook] = append(p.commands[hook], command{
			name: parts[0],
			args: parts[1:],
			mode: mode,
		})
	}

	return nil
}
