package main

import (
	"encoding/json"
	"fmt"
	"os"
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

func createAuthJsonAction(c *config) error {
	var reposWithAuth []repo
	for _, repo := range *c.starter.composer.repos {
		if repo.auth {
			reposWithAuth = append(reposWithAuth, repo)
		}
	}
	if len(reposWithAuth) < 1 {
		return nil
	}
	authConfig := AuthConfig{
		HTTPBasic: map[string]HTTPBasicCredentials{},
	}
	for _, repo := range reposWithAuth {
		var username, password string
		if err := huh.NewInput().
			Title(fmt.Sprintf("Enter your %s Username", repo.name)).
			Value(&username).
			Run(); err != nil {
			return err
		}
		if err := huh.NewInput().
			EchoMode(huh.EchoModePassword).
			Title(fmt.Sprintf("Enter your %s Password or License key", repo.name)).
			Value(&password).
			Run(); err != nil {
			return err
		}
		authConfig.HTTPBasic[repo.url] = HTTPBasicCredentials{
			Username: strings.TrimSpace(username),
			Password: strings.TrimSpace(password),
		}
	}
	authJSON, err := json.MarshalIndent(authConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling auth config: %w", err)
	}
	if err := os.WriteFile("auth.json", authJSON, 0644); err != nil {
		return fmt.Errorf("error writing auth.json: %w", err)
	}
	return nil
}
