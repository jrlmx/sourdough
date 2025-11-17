package main

import (
	"flag"
	"fmt"
	"sort"

	"github.com/jrlmx/sourdough/internal/cli"
)

func helpCommand(sd cli.SourdoughConfig) error {
	commands := getCliCommands()
	fmt.Println("Sourdough CLI:")
	flag.Usage()
	fmt.Println("Commands:")
	var sorted []cli.Command
	for _, cmd := range commands {
		sorted = append(sorted, cmd)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Name < sorted[j].Name
	})
	for _, cmd := range sorted {
		fmt.Printf("  %s: %s\n\n", cmd.Name, cmd.Desc)
	}
	return nil
}
