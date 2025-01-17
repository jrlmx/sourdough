package main

func handleUserCommands(p *project) error {
	return runUserCommands("default", p.commands)
}
