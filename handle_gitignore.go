package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func handleGitignore(p *project) error {
	path := filepath.Join(p.dir, ".gitignore")

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
