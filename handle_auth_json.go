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

		authConfig.HTTPBasic[repo.Url] = HTTPBasicCredentials{
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

	if err := ensureAuthJsonInGitignore(); err != nil {
		return err
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
		),
		huh.NewGroup(
			huh.NewInput().Title(fmt.Sprintf("Enter your %s Password or License key", repo.Name)).EchoMode(huh.EchoModePassword).Value(&password),
		),
	)

	if err := form.Run(); err != nil {
		return "", "", fmt.Errorf("failed receiving %s credentials: %w", repo.Name, err)
	}

	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}

func ensureAuthJsonInGitignore() error {
	path := filepath.Join(".", ".gitignore")

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read .gitignore: %w", err)
	}

	if strings.Contains(string(content), "auth.json") {
		fmt.Println("skipped: adding auth.json to .gitignore - already present")
		return nil
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open .gitignore: %w", err)
	}

	defer file.Close()
	if _, err := file.WriteString("\nauth.json\n"); err != nil {
		return fmt.Errorf("failed to update .gitignore: %w", err)
	}

	fmt.Println("Added auth.json to .gitignore")

	return nil
}
