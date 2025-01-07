package main

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh"
)

func handleKitSelection(p *project) error {
	opts, err := availableKits()

	if err != nil {
		return err
	}

	var kit string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a starter kit to apply").
				Options(huh.NewOptions(opts...)...).
				Value(&kit),
		),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("starter kit selection failed: %w", err)
	}
	p.kit = &kit
	p.loadConfig()

	return nil
}

func availableKits() ([]string, error) {
	opts := []string{}
	files, err := kits.ReadDir("kits")
	if err != nil {
		return []string{}, fmt.Errorf("error reading embedded directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			opts = append(opts, file.Name())
		}
	}

	if len(opts) < 1 {
		return []string{}, errors.New("no starter kits found")
	}

	return opts, nil
}
