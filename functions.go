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
