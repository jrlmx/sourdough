package starter

import (
	"context"
	"os"
	"os/exec"
)

type Command struct {
	name string
	args []string
}

func NewCommand(name string, args []string) (*Command, error) {
	return &Command{
		name: name,
		args: args,
	}, nil
}

func (c *Command) Run(ctx context.Context) error {
	return RunCommand(ctx, c.name, c.args)
}

func RunCommand(ctx context.Context, name string, args []string) error {
	c := exec.CommandContext(ctx, name, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}
	return nil
}

func RunCommandGroup(ctx context.Context, commands []Command) error {
	for _, cmd := range commands {
		if err := RunCommand(ctx, cmd.name, cmd.args); err != nil {
			return err
		}
	}
	return nil
}
