package main

import (
	"errors"
	"fmt"

	"github.com/jrlmx/sourdough/internal/cli"
)

func listCommand(sd cli.SourdoughConfig) error {
	starters := sd.StarterOptions()
	if len(starters) < 1 {
		return errors.New("no starters found")
	}
	for _, starter := range starters {
		fmt.Println(starter)
	}
	return nil
}
