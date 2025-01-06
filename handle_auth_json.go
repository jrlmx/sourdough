package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/charmbracelet/huh"
)

type AuthConfig struct {
	HTTPBasic map[string]HTTPBasicCredentials `json:"http-basic"`
}

type HTTPBasicCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func handleAuthJSON(p *project) error {
	if !slices.Contains(p.opts.PHP.Prod, "livewire/flux-pro") {
		fmt.Println("Skipping auth.json creation - Flux UI Pro is not in the project dependencies.")
		return nil
	}

	fmt.Println("Creating auth.json...")

	username, license, err := fluxPrompt()

	if err != nil {
		return err
	}

	authConfig := AuthConfig{
		HTTPBasic: map[string]HTTPBasicCredentials{
			"composer.fluxui.dev": {
				Username: username,
				Password: license,
			},
		},
	}

	path := filepath.Join(p.dir, "auth.json")

	file, err := json.MarshalIndent(authConfig, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to create auth.json: %w", err)
	}

	if err := os.WriteFile(path, file, 0644); err != nil {
		return fmt.Errorf("failed to write auth.json: %w", err)
	}

	return nil
}

func fluxPrompt() (string, string, error) {
	var (
		username string
		password string
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Enter your FLux UI Username").EchoMode(huh.EchoModePassword).Value(&username),
			huh.NewInput().Title("Enter your Flux Licence Key").EchoMode(huh.EchoModePassword).Value(&password),
		),
	)

	if err := form.Run(); err != nil {
		return "", "", fmt.Errorf("failed receiving flux credentials: %w", err)
	}

	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}
