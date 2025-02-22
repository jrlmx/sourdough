package main

import "fmt"

func npmInstallAction(cfg *config) error {
	fmt.Println("Removing unwanted npm packages...")
	if len(cfg.starter.npm.remove) > 0 {
		if err := runCommand("npm", append([]string{"remove", "--no-package-lock"}, cfg.starter.npm.remove...), QuietMode); err != nil {
			return fmt.Errorf("error removing unwanted npm packages: %w", err)
		}
	}
	fmt.Println("Installing npm packages...")
	if err := runCommand("npm", append([]string{"install", "--no-package-lock"}, cfg.starter.npm.production...), QuietMode); err != nil {
		return fmt.Errorf("error installing npm packages: %w", err)
	}
	if err := runCommand("npm", append([]string{"install", "-D", "--no-package-lock"}, cfg.starter.npm.development...), QuietMode); err != nil {
		return fmt.Errorf("error installing npm dev packages: %w", err)
	}
	if err := runCommand("npm", []string{"install"}, InteractiveMode); err != nil {
		return fmt.Errorf("error installing npm packages: %w", err)
	}
	return nil
}
