package main

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/jrlmx/sourdough/internal/cli"
	"github.com/jrlmx/sourdough/internal/starter"
)

func newCommand(sd cli.SourdoughConfig) error {
	options := sd.StarterOptions()
	if len(options) < 1 {
		return errors.New("no starters yet. try adding one with the 'add' command")
	}
	var sname, pname string
	if len(sd.Args) > 0 {
		sname = sd.Args[0]
	}
	if len(sd.Args) > 1 {
		pname = sd.Args[1]
	}
	if err := starterInput(options)(&sname); err != nil {
		return err
	}
	if err := cli.TextInput("project", []cli.Rule{
		cli.RequiredRule(),
		cli.BetweenRule(3, 255),
	})(&pname); err != nil {
		return err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	s, err := starter.NewStarter(
		pname,
		filepath.Join(sd.StarterPath(), sname),
		filepath.Join(pwd, pname),
	)
	if err != nil {
		return err
	}
	actions := []starter.Action{
		{
			Hook:     "create_project",
			Callback: starter.CreateNewProjectAction,
		},
		{
			Hook:     "navigate_to_project",
			Callback: starter.NavigateToProjectDirAction,
		},
		{
			Hook:     "remove_files",
			Callback: starter.RemoveFilesAction,
		},
		{
			Hook:     "php_dependencies",
			Callback: starter.PHPDependenciesAction,
		},
		{
			Hook:     "js_dependencies",
			Callback: starter.JSDependenciesAction,
		},
		{
			Hook:     "copy_files",
			Callback: starter.CopyFilesAction,
		},
		{
			Hook:     "run_commands",
			Callback: starter.RunCommandsAction,
		},
	}
	for _, action := range actions {
		if err := action.Callback(&sd, s); err != nil {
			return err
		}
	}
	return nil
}
