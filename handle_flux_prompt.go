package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func handleFluxPrompt(cfg *config) error {
	fmt.Print("Enter Flux UI Username (email): ")
	username, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to read username: %w", err)
	}
	fmt.Println()

	fmt.Print("Enter Flux UI License Key: ")
	license, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to read license key: %w", err)
	}
	fmt.Println()

	cfg.flux = &struct {
		username string
		password string
	}{
		username: strings.TrimSpace(string(username)),
		password: strings.TrimSpace(string(license)),
	}

	return nil
}
