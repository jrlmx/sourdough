package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type ValidationError struct {
	Name  string
	Value string
	Err   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("error validating %s: %s", e.Name, e.Err)
}

func Validate(key string, value string, rules ...Rule) error {
	return Validator(key, rules...)(value)
}

func Validator(key string, rules ...Rule) func(string) error {
	return func(value string) error {
		for _, rule := range rules {
			if err := rule(value); err != nil {
				return &ValidationError{
					Name:  key,
					Value: value,
					Err:   err.Error(),
				}
			}
		}
		return nil
	}
}

type Rule func(string) error

func BetweenRule(min, max int) Rule {
	return func(value string) error {
		if len(value) < min || len(value) > max {
			return fmt.Errorf("must be between %d and %d", min, max)
		}
		return nil
	}
}

func RequiredRule() Rule {
	return func(value string) error {
		if value == "" {
			return fmt.Errorf("cannot be empty")
		}
		return nil
	}
}

type PathOption int

const (
	Any PathOption = iota
	Absolute
	Relative
)

func PathRule(opt PathOption) Rule {
	return func(value string) error {
		cleaned := filepath.Clean(value)
		if opt == Absolute {
			if !filepath.IsAbs(cleaned) {
				return fmt.Errorf("must be an absolute path")
			}
		}
		if opt == Relative {
			if filepath.IsAbs(cleaned) {
				return fmt.Errorf("must be a relative path")
			}
		}
		return nil
	}
}

func PathExistsRule() Rule {
	return func(value string) error {
		if _, err := os.Stat(value); os.IsNotExist(err) {
			return fmt.Errorf("path does not exist")
		}
		return nil
	}
}

func PathNotExistsRule() Rule {
	return func(value string) error {
		if _, err := os.Stat(value); err == nil {
			return fmt.Errorf("path already exists")
		}
		return nil
	}
}

func InRule(values ...string) Rule {
	return func(value string) error {
		if slices.Contains(values, value) {
			return nil
		}
		return fmt.Errorf("must be one of %v", values)
	}
}

func NotInRule(values ...string) Rule {
	return func(value string) error {
		if !slices.Contains(values, value) {
			return nil
		}
		return fmt.Errorf("must not be one of %v", values)
	}
}

func IntRule(min, max int) Rule {
	return func(value string) error {
		if len(value) < min || len(value) > max {
			return fmt.Errorf("must be between %d and %d", min, max)
		}
		return nil
	}
}

func GitRepoRule() Rule {
	return func(value string) error {
		if value == "" {
			return fmt.Errorf("git repo URL cannot be empty")
		}

		if strings.HasPrefix(value, "https://") || strings.HasPrefix(value, "http://") {
			parts := strings.Split(value, "://")
			if len(parts) != 2 {
				return fmt.Errorf("invalid HTTPS git URL")
			}
			remainder := parts[1]
			if !strings.Contains(remainder, "/") {
				return fmt.Errorf("invalid HTTPS git URL")
			}
			pathParts := strings.Split(remainder, "/")
			if len(pathParts) < 3 {
				return fmt.Errorf("invalid HTTPS git URL")
			}
			repo := strings.TrimSuffix(pathParts[len(pathParts)-1], ".git")
			if repo == "" {
				return fmt.Errorf("invalid HTTPS git URL")
			}
			return nil
		}

		if strings.HasPrefix(value, "git@") {
			parts := strings.Split(value, ":")
			if len(parts) != 2 {
				return fmt.Errorf("invalid SSH git URL")
			}
			hostPart := parts[0]
			pathPart := parts[1]
			if !strings.Contains(hostPart, "@") || !strings.Contains(pathPart, "/") {
				return fmt.Errorf("invalid SSH git URL")
			}
			pathParts := strings.Split(pathPart, "/")
			if len(pathParts) < 2 {
				return fmt.Errorf("invalid SSH git URL")
			}
			repo := strings.TrimSuffix(pathParts[len(pathParts)-1], ".git")
			if repo == "" {
				return fmt.Errorf("invalid SSH git URL")
			}
			return nil
		}

		if filepath.IsAbs(value) || strings.HasPrefix(value, ".") {
			if !strings.Contains(value, "://") {
				return nil
			}
		}

		return fmt.Errorf("invalid git repository URL format")
	}
}
