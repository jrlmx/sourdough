package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func handlePublishFiles(cfg *config) error {
	err := os.MkdirAll(cfg.projectDir, os.ModePerm)

	if err != nil {
		return fmt.Errorf("failed to create desination directory: %w", err)
	}

	err = fs.WalkDir(stubs, "stubs", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking directory: %w", err)
		}

		if d.IsDir() {
			return nil
		}

		content, err := stubs.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", path, err)
		}

		relPath, err := filepath.Rel("stubs", path)
		if err != nil {
			return fmt.Errorf("failed to determine relative path: %w", err)
		}
		destPath := filepath.Join(cfg.projectDir, relPath)

		err = os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", destPath, err)
		}

		err = os.WriteFile(destPath, content, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", destPath, err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	fmt.Println("Files copied successfully.")

	return nil
}
