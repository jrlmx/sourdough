package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func handleAuthJSON(cfg *config) error {
	if slices.ContainsFunc(cfg.deps.Composer, func(s string) bool {
		return strings.Contains(s, "livewire/flux")
	}) {
		return nil
	}

	authConfig := AuthConfig{
		HTTPBasic: map[string]HTTPBasicCredentials{
			"composer.fluxui.dev": {
				Username: cfg.flux.username,
				Password: cfg.flux.password,
			},
		},
	}

	path := filepath.Join(cfg.projectDir, "auth.json")

	file, err := json.MarshalIndent(authConfig, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to create auth.json: %w", err)
	}

	if err := os.WriteFile(path, file, 0644); err != nil {
		return fmt.Errorf("failed to write auth.json: %w", err)
	}

	return nil
}
