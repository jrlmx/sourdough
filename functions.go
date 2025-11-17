package main

import "github.com/jrlmx/sourdough/internal/cli"

func starterInput(options []string) cli.Input {
	return cli.SelectInput("starter name", options, []cli.Rule{
		cli.RequiredRule(),
		cli.InRule(options...),
	})
}
