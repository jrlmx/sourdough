package main

import (
	"errors"
	"fmt"
	"slices"

	"github.com/charmbracelet/huh"
)

func handleStarterSelection(p *project) error {
	opts, err := availableStarters()

	if err != nil {
		return err
	}

	if len(opts) < 1 {
		return errors.New("no starters found")
	}

	kit, err := promptForStarterKit(opts)
	if err != nil {
		return err
	}

	p.init(kit)

	return nil
}

func availableStarters() ([]string, error) {
	opts := []string{}
	files, err := starters.ReadDir("starters")
	if err != nil {
		return []string{}, fmt.Errorf("error reading embedded directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			opts = append(opts, file.Name())
		}
	}

	if len(opts) < 1 {
		return []string{}, errors.New("no starters found")
	}

	return opts, nil
}

func promptForStarterKit(opts []string) (string, error) {
	var kit string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a starter kit to apply").
				Options(huh.NewOptions(opts...)...).
				Value(&kit).Validate(func(s string) error {
				if !slices.Contains(opts, s) {
					return fmt.Errorf("invalid starter kit: %s", s)
				}
				return nil
			}),
		),
	)

	if err := form.Run(); err != nil {
		return "", fmt.Errorf("starter kit selection failed: %w", err)
	}

	return kit, nil
}
