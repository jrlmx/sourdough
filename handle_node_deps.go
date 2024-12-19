package main

import (
	"fmt"
	"os"
	"os/exec"
)

func handleNodeDeps(cfg *config) error {
	if err := exec.Command("which", "npm").Run(); err != nil {
		return fmt.Errorf("npm is not installed")
	}

	deps := []string{
		"prettier",
		"prettier-plugin-blade",
		"@tailwindcss/typography",
		"@tailwindcss/forms",
	}

	cmd := exec.Command("npm", append([]string{"install", "-D"}, deps...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install node dependencies: %w", err)
	}

	return nil
}
