package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jrlmx/sourdough/internal/cli"
)

func addCommand(sd cli.SourdoughConfig) error {
	fmt.Println("Add a new Starter")
	var sname, repoUrl string
	if len(sd.Args) > 0 {
		sname = sd.Args[0]
	}
	if len(sd.Args) > 1 {
		repoUrl = sd.Args[1]
	}
	if err := cli.TextInput("starter name", []cli.Rule{
		cli.RequiredRule(),
		cli.BetweenRule(3, 255),
		cli.NotInRule(sd.StarterOptions()...),
	})(&sname); err != nil {
		return err
	}
	if err := cli.TextInput("repository", []cli.Rule{
		cli.RequiredRule(),
		cli.GitRepoRule(),
	})(&repoUrl); err != nil {
		return err
	}
	spath := filepath.Join(sd.StarterPath(), sname)
	sd.CM.Add("delete_starter_folder", func() error {
		return os.RemoveAll(spath)
	})
	if err := exec.Command("git", "clone", "--depth", "1", repoUrl, spath).Run(); err != nil {
		return err
	}
	if _, err := os.Stat(filepath.Join(spath, "starter.json")); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("starter '%s' does not contain a starter.json file", sname)
		}
		return err
	}
	sd.CM.Remove("delete_starter_folder")
	fmt.Println("Starter added successfully")
	return nil
}
