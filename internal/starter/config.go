package starter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type PackageManifest struct {
	production  []string
	development []string
	remove      []string
}

type FileManifest struct {
	remove []string
}

type StarterConfig struct {
	target   string
	source   string
	install  Command
	files    FileManifest
	php      PackageManifest
	js       PackageManifest
	stubs    []string
	commands map[string][]Command
}

type StarterConfigJson struct {
	Installer struct {
		Args []string `json:"args"`
	} `json:"installer"`
	Files struct {
		Remove []string `json:"remove"`
	} `json:"files"`
	PHP struct {
		Production  []string `json:"production"`
		Development []string `json:"development"`
		Remove      []string `json:"remove"`
	} `json:"php"`
	JS struct {
		Production  []string `json:"production"`
		Development []string `json:"development"`
		Remove      []string `json:"remove"`
	} `json:"js"`
	Commands []string `json:"commands"`
}

func NewStarter(spath string, target string) (*StarterConfig, error) {
	sjson, err := getStarterConfig(filepath.Join(spath, "starter.json"))
	if err != nil {
		return nil, err
	}

	sfiles, err := getStarterStubs(filepath.Join(spath, "stubs"))
	if err != nil {
		return nil, err
	}

	installCmd, err := NewCommand("laravel", append([]string{"new"}, sjson.Installer.Args...))
	if err != nil {
		return nil, err
	}

	commands := map[string][]Command{}
	for _, cstring := range sjson.Commands {
		trimmed := strings.TrimSpace(cstring)
		if trimmed == "" {
			continue
		}
		hook := "default"
		parts := strings.Split(cstring, " ")
		if len(parts) > 1 {
			if strings.HasPrefix(parts[0], "@") {
				hook = strings.TrimLeft(parts[0], "@")
				parts = parts[1:]
			}
		}
		cname := parts[0]
		var cargs []string
		if len(parts) > 1 {
			cargs = parts[1:]
		}
		command, err := NewCommand(cname, cargs)
		if err != nil {
			return nil, err
		}
		if _, ok := commands[hook]; !ok {
			commands[hook] = []Command{}
		}
		commands[hook] = append(commands[hook], *command)
	}

	return &StarterConfig{
		target:  target,
		source:  spath,
		install: *installCmd,
		files: FileManifest{
			remove: sjson.Files.Remove,
		},
		php: PackageManifest{
			production:  sjson.PHP.Production,
			development: sjson.PHP.Development,
			remove:      sjson.PHP.Remove,
		},
		js: PackageManifest{
			production:  sjson.JS.Production,
			development: sjson.JS.Development,
			remove:      sjson.JS.Remove,
		},
		stubs:    sfiles,
		commands: commands,
	}, nil
}

func getStarterConfig(filename string) (*StarterConfigJson, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errors.New("missing starter.json file")
		}
		return nil, err
	}
	var sjson StarterConfigJson
	err = json.Unmarshal(file, &sjson)
	if err != nil {
		return nil, err
	}
	return &sjson, nil
}

func getStarterStubs(dirname string) ([]string, error) {
	sfiles := []string{}
	if err := filepath.WalkDir(dirname, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		sfiles = append(sfiles, path)
		return nil
	}); err != nil {
		return nil, fmt.Errorf("error getting starter files: %s", err)
	}
	return sfiles, nil
}
