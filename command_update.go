package main

import (
	"errors"

	"github.com/jrlmx/sourdough/internal/cli"
)

func updateCommand(sd cli.SourdoughConfig) error {
	options := sd.StarterOptions()
	if len(options) < 1 {
		return errors.New("no starters yet. try adding one with the 'add' command")
	}
	var sname string
	if len(sd.Args) > 0 {
		sname = sd.Args[0]
	}
	if err := starterInput(options)(&sname); err != nil {
		return err
	}
	// Fetch most recent changes from the repo...
	return nil
}
