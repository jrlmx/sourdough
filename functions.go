package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func run(cname string, cargs ...string) error {
	cmd := exec.Command(cname, cargs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func runQuietly(cname string, cargs ...string) error {
	cmd := exec.Command(cname, cargs...)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func runInteractive(cname string, cargs ...string) error {
	cmd := exec.Command(cname, cargs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func check(cname string) (string, error) {
	path, err := exec.LookPath(cname)
	if err != nil {
		return "", fmt.Errorf("%s is not installed", cname)
	}

	return path, nil
}

func cleanString(s string) string {
	s = strings.TrimSpace(s)

	if s != "" {
		re := regexp.MustCompile(`[<>:"/\|?*]`)
		s = re.ReplaceAllString(s, "_")
	}

	return s
}

func existsAndNotEmpty(dir string) (bool, error) {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if !info.IsDir() {
		return false, fmt.Errorf("%s is not a directory", dir)
	}

	contents, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}

	return len(contents) > 0, nil
}
