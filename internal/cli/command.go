package cli

type Command struct {
	Name    string
	Desc    string
	handler HandlerFunc
}

type HandlerFunc func(config SourdoughConfig) error

func NewCommand(name, desc string, handler HandlerFunc) *Command {
	return &Command{
		Name:    name,
		Desc:    desc,
		handler: handler,
	}
}

func (c *Command) Exec(cfg *SourdoughConfig) error {
	return c.handler(*cfg)
}
