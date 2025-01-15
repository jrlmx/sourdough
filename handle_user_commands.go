package main

func handleUserCommands(p *project) error {
	for _, cmd := range p.config.Commands {
		if err := runUntrustedCommand(cmd); err != nil {
			return err
		}
	}

	return nil
}
