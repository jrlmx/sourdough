package main

import "fmt"

func runCommandsAction(c *config) error {
	for _, cmd := range c.starter.commands["default"] {
		if err := cmd.run(); err != nil {
			return fmt.Errorf("error running command: %w", err)
		}
	}
	return nil
}
