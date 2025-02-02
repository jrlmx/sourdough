package main

import (
	"flag"
	"fmt"
)

type cliOptions struct {
	projectName    string
	starterName    string
	printHooks     bool
	addStarter     bool
	removeStarter  string
	exportStarters bool
	importStarters string
}

func getCliOptions() cliOptions {
	var projectName string
	var starterName string
	var printHooks bool
	var addStarter bool
	var removeStarter string
	var exportStarters bool
	var importStarters string

	// Options
	flag.StringVar(&starterName, "starter", "", "Name of the starter")
	flag.BoolVar(&printHooks, "hooks", false, "Print available hooks")

	// Sub commands
	flag.BoolVar(&addStarter, "add", false, "Add a new starter")
	flag.StringVar(&removeStarter, "remove", "", "Remove a starter")
	flag.BoolVar(&exportStarters, "export", false, "Export starters to sourdough.json")
	flag.StringVar(&importStarters, "import", "", "Import starters from sourdough.json")

	flag.Usage = func() {
		fmt.Println("Usage: sourdough [options] [project-name]")

		fmt.Println("\nOptions:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(flag.Args()) > 0 {
		projectName = flag.Args()[0]
	}

	return cliOptions{
		projectName,
		starterName,
		printHooks,
		addStarter,
		removeStarter,
		exportStarters,
		importStarters,
	}
}
