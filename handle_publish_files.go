package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func handlePublishFiles(p *project) error {
	dir := filepath.Join(".")
	stubs := filepath.Join("kits", *p.kit, "stubs")

	err := fs.WalkDir(kits, stubs, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking directory %s: %w", path, err)
		}

		if d.IsDir() {
			return nil
		}

		content, err := kits.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", path, err)
		}

		relPath, err := filepath.Rel(stubs, path)
		if err != nil {
			return fmt.Errorf("failed to determine relative path: %w", err)
		}

		destPath := filepath.Join(dir, relPath)
		if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", destPath, err)
		}

		if err := os.WriteFile(destPath, content, os.ModePerm); err != nil {
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
