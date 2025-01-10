package main

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	PHP     PackageManifest `json:"php"`
	JS      PackageManifest `json:"js"`
	Files   []string        `json:"remove_files"`
	Repos   []Repo          `json:"repos"`
	Artisan []string        `json:"artisan"`
	NPX     []string        `json:"npx"`
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
	pdir   string
	kit    *string
	config *Config
}

func newProject(name, pdir string) project {
	return project{
		name,
		pdir,
		nil,
		nil,
	}
}

func (p *project) loadConfig() error {
	data, err := starters.ReadFile("starters/" + *p.kit + "/config.json")
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
