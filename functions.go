package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"
)

type runMode int

const (
	NormalMode runMode = iota
	QuietMode
	InteractiveMode
)

func runCommand(cname string, mode runMode, cargs ...string) error {
	cmd := exec.Command(cname, cargs...)
	if SHOUT || mode != QuietMode {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if mode == InteractiveMode {
		cmd.Stdin = os.Stdin
	}
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func runUserCommand(cmd command) error {
	err := validateCommand(cmd.name, cmd.args)
	if err != nil {
		return fmt.Errorf("invalid command: %w", err)
	}

	return runCommand(cmd.name, cmd.mode, cmd.args...)
}

func runUserCommands(hook string, commands map[string][]command) error {
	if _, ok := commands[hook]; !ok {
		return nil
	}

	if len(commands[hook]) < 1 {
		return nil
	}

	for _, cmd := range commands[hook] {
		if err := runUserCommand(cmd); err != nil {
			return err
		}
	}

	return nil
}

func validateCommand(cname string, cargs []string) error {
	if len(cname) < 1 {
		return errors.New("command name cannot be empty")
	}

	if strings.Contains(cname, "..") || strings.HasPrefix(cname, "/") {
		return fmt.Errorf("command attempted to ascend to parent directory: %s", cname)
	}

	if strings.ContainsAny(strings.Join(append([]string{cname}, cargs...), ""), "&;|<>(){}$`\\") {
		return fmt.Errorf("command attempted to ascend to parent directory: %s", cname)
	}

	allowed := []string{
		"php",
		"composer",
		"npm",
		"npx",
		"git",
	}

	allowedPrefixes := []string{
		"./vendor/",
		"./node_modules/",
	}

	if !slices.Contains(allowed, cname) && !slices.ContainsFunc(allowedPrefixes, func(s string) bool {
		return strings.HasPrefix(cname, s)
	}) {
		return fmt.Errorf("command not allowed: %s", cname)
	}

	if slices.Contains(allowed, cname) {
		if _, err := isInstalled(cname); err != nil {
			return err
		}
		return nil
	}

	if slices.ContainsFunc(allowedPrefixes, func(s string) bool {
		return strings.HasPrefix(cname, s)
	}) {
		if _, err := os.Stat(cname); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("command %s is not installed", cname)
			}
			return fmt.Errorf("error checking command %s: %w", cname, err)
		}
		return nil
	}

	return errors.New("command not allowed")
}

func isInstalled(cname string) (string, error) {
	path, err := exec.LookPath(cname)
	if err != nil {
		return "", fmt.Errorf("%s is not installed", cname)
	}

	return path, nil
}

func cleanDirName(s string) string {
	s = strings.TrimSpace(s)

	if s != "" {
		re := regexp.MustCompile(`[<>:"/\|?*]`)
		s = re.ReplaceAllString(s, "_")
	}

	return s
}
