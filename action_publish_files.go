package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func publishFilesAction(cfg *config) error {
	stubs := filepath.Join(cfg.starter.dir, "stubs")
	info, err := os.Stat(stubs)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("No stubs directory found. Skipping...")
			return nil
		}
		return fmt.Errorf("error accessing stubs directory: %w", err)
	}
	if !info.IsDir() {
		fmt.Println("Stubs path is not a directory. Skipping...")
		return nil
	}
	if err = fs.WalkDir(os.DirFS(stubs), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking directory: %w", err)
		}

		if d.IsDir() {
			return nil
		}

		src := filepath.Join(stubs, path)
		dest := filepath.Join(cfg.project.dir, path)
		if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
			return fmt.Errorf("error creating directory: %w", err)
		}

		srcFile, err := os.ReadFile(src)
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}

		if err := os.WriteFile(dest, srcFile, os.ModePerm); err != nil {
			return fmt.Errorf("error writing file: %w", err)
		}

		return nil
	}); err != nil {
		return err
	}
	fmt.Println("Files copied successfully.")
	return nil
}
