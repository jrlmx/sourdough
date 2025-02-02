package main

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

type CommandMode int

const (
	QuietMode CommandMode = iota
	NormalMode
	InteractiveMode
)

type command struct {
	name string
	args []string
	mode CommandMode
}

func newCommand(name string, args []string, mode CommandMode) *command {
	cmd := command{
		name: name,
		args: args,
		mode: mode,
	}
	return &cmd
}

func runCommand(name string, args []string, mode CommandMode) error {
	cmd := newCommand(name, args, mode)
	if err := cmd.validate(); err != nil {
		return err
	}
	return cmd.run()
}

func (c *command) validate() error {
	if strings.Contains(c.name, "..") {
		return fmt.Errorf("command should not ascend directories: %s", c.name)
	}
	if strings.ContainsAny(c.name+strings.Join(c.args, ""), "&;|<>(){}$`") {
		return fmt.Errorf("command contains invalid characters: %s [%s]", c.name, c.args)
	}
	if !slices.ContainsFunc([]string{
		"php",
		"composer",
		"laravel",
		"npm",
		"npx",
		"git",
		"./vendor/",
		"./node_modules/",
	}, func(s string) bool {
		if c.name == s || strings.HasPrefix(c.name, s) {
			return true
		}
		return false
	}) {
		return fmt.Errorf("command not allowed: %s", c.name)
	}
	return nil
}

func (c *command) run() error {
	cmd := exec.Command(c.name, c.args...)
	cmd.Stderr = os.Stderr
	if c.mode != QuietMode {
		cmd.Stdout = os.Stdout
	}
	if c.mode == InteractiveMode {
		cmd.Stdin = os.Stdin
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error executing command %s: %w", c.name, err)
	}
	return nil
}
