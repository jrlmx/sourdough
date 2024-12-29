package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func handleAuthJSON(cfg *config) error {
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

func fluxPrompt() (string, string, error) {
	r := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Flux UI Username (email): ")
	username, _ := r.ReadString('\n')

	fmt.Print("Enter Flux UI License Key: ")
	license, _ := r.ReadString('\n')

	return strings.TrimSpace(username), strings.TrimSpace(license), nil
}
