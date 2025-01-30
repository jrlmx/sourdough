package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func removeFilesAction(cfg *config) error {
	fmt.Println("Removing unwanted files...")
	for _, file := range cfg.starter.remove {
		cleanPath := filepath.Join(".", file)
		if err := os.RemoveAll(cleanPath); err != nil {
			return fmt.Errorf("error removing file: %w", err)
		}
	}
	return nil
}
