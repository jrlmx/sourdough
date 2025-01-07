package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
	var repos []Repo

	for _, repo := range p.config.Repos {
		if repo.Auth {
			repos = append(repos, repo)
		}
	}

	if len(repos) < 1 {
		fmt.Println("No repositories need authorization")
		return nil
	}

	authConfig := AuthConfig{
		HTTPBasic: map[string]HTTPBasicCredentials{},
	}

	for _, repo := range repos {
		username, password, err := authPrompt(repo)

		if err != nil {
			return err
		}

		authConfig.HTTPBasic[repo.Name] = HTTPBasicCredentials{
			Username: username,
			Password: password,
		}
	}

	fmt.Println("Creating auth.json...")

	path := filepath.Join(".", "auth.json")
	file, err := json.MarshalIndent(authConfig, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to create auth.json: %w", err)
	}

	if err := os.WriteFile(path, file, 0644); err != nil {
		return fmt.Errorf("failed to write auth.json: %w", err)
	}

	return nil
}

func authPrompt(repo Repo) (string, string, error) {
	var (
		username string
		password string
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title(fmt.Sprintf("Enter your %s Username", repo.Name)).EchoMode(huh.EchoModePassword).Value(&username),
			huh.NewInput().Title(fmt.Sprintf("Enter your %s Password", repo.Name)).EchoMode(huh.EchoModePassword).Value(&password),
		).Title("Enter your " + repo.Name + " credentials"),
	)

	if err := form.Run(); err != nil {
		return "", "", fmt.Errorf("failed receiving %s credentials: %w", repo.Name, err)
	}

	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}
