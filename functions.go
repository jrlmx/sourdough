package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
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

func check(cname string) (string, error) {
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
