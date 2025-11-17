package main

import (
	"context"
	"log"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/jrlmx/sourdough/internal/cli"
)

func main() {
	if err := checkDependencies(); err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	commands := getCliCommands()
	cfg := cli.NewSourdoughConfig(ctx)

	defer func() {
		if err := cfg.CM.Clean(); err != nil {
			log.Fatal(err)
		}
	}()

	cname := cfg.Cmd
	if cname == "" {
		cname = "help"
	}
	cmd, ok := commands[cname]
	if !ok {
		log.Fatalf("command %s not found", cname)
	}
	if err := cmd.Exec(cfg); err != nil {
		log.Fatal(err)
	}

	stop()
	<-ctx.Done()
}

func checkDependencies() error {
	dependencies := []string{
		"laravel",
		"composer",
		"npm",
		"git",
	}
	for _, dependency := range dependencies {
		if _, err := exec.LookPath(dependency); err != nil {
			return err
		}
	}
	return nil
}

func getCliCommands() map[string]cli.Command {
	return map[string]cli.Command{
		"install": *cli.NewCommand(
			"install",
			"initialize sourdough's local data",
			installCommand,
		),
		"help": *cli.NewCommand(
			"help",
			"display this help message -> usage: help",
			helpCommand,
		),
		"list": *cli.NewCommand(
			"list",
			"list available starters -> usage: list",
			listCommand,
		),
		"new": *cli.NewCommand(
			"new",
			"create a new project. -> usage: new <starter_name> <project_name>",
			newCommand,
		),
		"apply": *cli.NewCommand(
			"apply",
			"apply a starter to an existing project. -> usage: apply <starter_name> <path>",
			applyCommand,
		),
		"add": *cli.NewCommand(
			"add",
			"add a starter from a git repository. -> usage: add <starter_name> <repository_url>",
			addCommand,
		),
		"remove": *cli.NewCommand(
			"remove",
			"remove a starter from disk. -> usage: remove <starter_name>",
			removeCommand,
		),
		"update": *cli.NewCommand(
			"update",
			"update a starter on disk. -> usage: update <starter_name> <?repository_url>",
			updateCommand,
		),
	}
}
