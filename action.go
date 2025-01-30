package main

type Action struct {
	name     string
	hookable bool
	callback func(c *config) error
}
