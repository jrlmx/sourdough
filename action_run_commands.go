package main

import "fmt"

func runCommandsAction(cfg *config) error {
	for _, cmd := range cfg.starter.commands["default"] {
		if err := cmd.run(); err != nil {
			return fmt.Errorf("error running command: %w", err)
		}
	}
	return nil
}
