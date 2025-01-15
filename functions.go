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
	if mode != QuietMode {
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

func run(cname string, cargs ...string) error {
	return runCommand(cname, NormalMode, cargs...)
}

func runQuietly(cname string, cargs ...string) error {
	return runCommand(cname, QuietMode, cargs...)
}

func runInteractive(cname string, cargs ...string) error {
	return runCommand(cname, InteractiveMode, cargs...)
}

func runUntrustedCommand(cmd string) error {
	parts := strings.Split(strings.TrimSpace(cmd), " ")
	cname := parts[0]
	cargs := parts[1:]

	var mode runMode

	if strings.HasPrefix(cname, "quietly:") {
		mode = QuietMode
		cname = strings.TrimPrefix(cname, "quietly:")
	} else if strings.HasPrefix(cname, "interactive:") {
		mode = InteractiveMode
		cname = strings.TrimPrefix(cname, "interactive:")
	} else {
		mode = NormalMode
	}

	err := validateCommand(cname, cargs)
	if err != nil {
		return err
	}

	return runCommand(cname, mode, cargs...)
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
		"laravel",
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
		return fmt.Errorf("invalid command: %s", cname)
	}

	if _, err := os.Stat(cname); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("command %s does not exist", cname)
		}
		return fmt.Errorf("error checking command %s: %w", cname, err)
	}

	return nil
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
