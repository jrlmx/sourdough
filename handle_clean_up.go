package main

import (
	"fmt"
	"os"
)

func handleCleanUp(cfg *config) error {
	fmt.Println("Cleaning up...")

	err := os.Remove(cfg.projectDir + "/resources/views/welcome.blade.php")

	if err != nil {
		return err
	}

	return nil
}
