package cli

import "github.com/charmbracelet/huh"

type Input func(*string) error

func TextInput(name string, rules []Rule) Input {
	return func(pointer *string) error {
		if err := Validate(name, *pointer, rules...); err != nil {
			if err := huh.NewInput().
				Value(pointer).
				Title(name).
				Validate(Validator(name, rules...)).
				Run(); err != nil {
				return err
			}
		}
		return nil
	}
}

func SelectInput(name string, options []string, rules []Rule) Input {
	return func(pointer *string) error {
		if err := Validate(name, *pointer, rules...); err != nil {
			var opts []huh.Option[string]
			for _, opt := range options {
				opts = append(opts, huh.Option[string]{
					Key:   opt,
					Value: opt,
				})
			}
			if err := huh.NewSelect[string]().
				Value(pointer).
				Title(name).
				Options(opts...).
				Validate(Validator(name, rules...)).
				Run(); err != nil {
				return err
			}
		}
		return nil
	}
}
