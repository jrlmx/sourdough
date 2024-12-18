package main

import (
	"embed"
	"fmt"
	"log"
	"os"
)

//go:embed stubs/*
var stubs embed.FS

type AuthConfig struct {
	HTTPBasic map[string]HTTPBasicCredentials `json:"http-basic"`
}

type HTTPBasicCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type config struct {
	dir  string
	flux *struct {
		username string
		password string
	}
}

func main() {
	dir, err := os.Getwd()

	if err != nil {
		log.Fatalf("Failed to get working director: %v", err)
	}

	cfg := config{
		dir: dir,
	}

	for _, action := range getActions() {
		err = action(&cfg)

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Laravel project scaffolding complete!")
}

func getActions() []func(cfg *config) error {
	return []func(cfg *config) error{
		isLaravelProject,
		fluxPrompt,
		createAuthJSON,
		updateGitignore,
		installComposerDependencies,
		copyStubFiles,
	}
}
