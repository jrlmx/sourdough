package main

import "fmt"

func composerInstallAction(cfg *config) error {
	fmt.Println("Adding composer repositories...")
	for _, repo := range *cfg.starter.composer.repos {
		if err := runCommand("composer", []string{"config", "repositories." + repo.name, "composer", "https://" + repo.url}, NormalMode); err != nil {
			return fmt.Errorf("error adding composer repository: %w", err)
		}
	}
	fmt.Println("Removing unwanted composer packages...")
	if len(cfg.starter.composer.remove) > 0 {
		if err := runCommand("composer", append([]string{"remove", "-n", "--no-update"}, cfg.starter.composer.remove...), QuietMode); err != nil {
			return fmt.Errorf("error removing unwanted composer packages: %w", err)
		}
	}
	fmt.Println("Installing composer packages...")
	if err := runCommand("composer", append([]string{"require", "-n", "--no-update"}, cfg.starter.composer.production...), QuietMode); err != nil {
		return fmt.Errorf("error requiring composer packages: %w", err)
	}
	if err := runCommand("composer", append([]string{"require", "-n", "--no-update", "--dev"}, cfg.starter.composer.development...), QuietMode); err != nil {
		return fmt.Errorf("error requiring composer dev packages: %w", err)
	}
	if err := runCommand("composer", []string{"update", "-n", "--no-scripts"}, QuietMode); err != nil {
		return fmt.Errorf("error updating composer packages: %w", err)
	}
	if err := runCommand("composer", []string{"install"}, NormalMode); err != nil {
		return fmt.Errorf("error installing composer packages: %w", err)
	}
	return nil
}
