package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func handleFluxPrompt(cfg *config) error {
	r := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Flux UI Username (email): ")
	username, _ := r.ReadString('\n')

	fmt.Print("Enter Flux UI License Key: ")
	license, _ := r.ReadString('\n')

	cfg.flux = &struct {
		username string
		password string
	}{
		username: strings.TrimSpace(username),
		password: strings.TrimSpace(license),
	}

	return nil
}
