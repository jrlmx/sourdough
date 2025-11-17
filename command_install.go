package main

import (
	"fmt"
	"os"

	"github.com/jrlmx/sourdough/internal/cli"
)

func installCommand(sd cli.SourdoughConfig) error {
	fmt.Println("Installing Sourdough's data folder(s) at ~/.sourdough and ~/.sourdough/starters")
	if sd.Flags.Force {
		fmt.Println("The force flag was used, removing ~/.sourdough if it exists")
		if err := os.RemoveAll(sd.DataPath()); err != nil {
			return err
		}
	}
	if err := os.MkdirAll(sd.StarterPath(), 0755); err != nil {
		return err
	}
	fmt.Println("Done.")
	return nil
}
