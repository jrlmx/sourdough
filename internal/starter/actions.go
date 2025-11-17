package starter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jrlmx/sourdough/internal/cli"
)

type Action struct {
	Hook     string
	Callback func(*cli.SourdoughConfig, *StarterConfig) error
}

func (a *Action) Run(sd *cli.SourdoughConfig, s *StarterConfig) error {
	if err := a.Callback(sd, s); err != nil {
		return err
	}
	if cgroup, ok := s.commands[a.Hook]; ok {
		if err := RunCommandGroup(sd.Ctx, cgroup); err != nil {
			return err
		}
	}
	return nil
}

func CreateNewProjectAction(sd *cli.SourdoughConfig, s *StarterConfig) error {
	fmt.Println("creating new project...")
	if err := s.install.Run(sd.Ctx); err != nil {
		return err
	}
	sd.CM.Add("delete_project_folder", func() error {
		fmt.Println("cleaning up project folder...")
		return os.RemoveAll(s.target)
	})
	return nil
}

func NavigateToProjectDirAction(sd *cli.SourdoughConfig, s *StarterConfig) error {
	fmt.Println("navigating to project directory...")
	if err := os.Chdir(s.target); err != nil {
		return nil
	}
	return nil
}

func RemoveFilesAction(sd *cli.SourdoughConfig, s *StarterConfig) error {
	fmt.Println("removing unwanted files...")
	for _, file := range s.files.remove {
		if strings.Contains(file, "..") {
			return fmt.Errorf("invalid file path '%s'", file)
		}
		cleaned := filepath.Join(".", file)
		if err := os.RemoveAll(cleaned); err != nil {
			return err
		}
	}
	return nil
}

func PHPDependenciesAction(sd *cli.SourdoughConfig, s *StarterConfig) error {
	fmt.Println("php dependencies...")
	if len(s.php.remove) > 0 {
		cstring := "remove -n --no-update " + strings.Join(s.php.remove, " ")
		if err := RunCommand(sd.Ctx, "composer", strings.Split(cstring, " ")); err != nil {
			return err
		}
	}
	if len(s.php.development) > 0 {
		cstring := "require -n --no-update --dev " + strings.Join(s.php.development, " ")
		if err := RunCommand(sd.Ctx, "composer", strings.Split(cstring, " ")); err != nil {
			return err
		}
	}
	if len(s.php.production) > 0 {
		cstring := "require -n --no-update " + strings.Join(s.php.production, " ")
		if err := RunCommand(sd.Ctx, "composer", strings.Split(cstring, " ")); err != nil {
			return err
		}
	}
	if err := RunCommand(sd.Ctx, "composer", []string{"install"}); err != nil {
		return err
	}
	return nil
}

func JSDependenciesAction(sd *cli.SourdoughConfig, s *StarterConfig) error {
	fmt.Println("js dependencies...")
	if len(s.js.remove) > 0 {
		cstring := "uninstall --no-package-lock " + strings.Join(s.js.remove, " ")
		if err := RunCommand(sd.Ctx, "npm", strings.Split(cstring, " ")); err != nil {
			return err
		}
	}
	if len(s.js.development) > 0 {
		cstring := "install --no-package-lock " + strings.Join(s.js.development, " ")
		if err := RunCommand(sd.Ctx, "npm", strings.Split(cstring, " ")); err != nil {
			return err
		}
	}
	if len(s.js.production) > 0 {
		cstring := "install --no-package-lock --save-dev "
		if err := RunCommand(sd.Ctx, "npm", strings.Split(cstring, " ")); err != nil {
			return err
		}
	}
	if err := RunCommand(sd.Ctx, "npm", []string{"install"}); err != nil {
		return err
	}
	return nil
}

func CopyFilesAction(sd *cli.SourdoughConfig, s *StarterConfig) error {
	fmt.Println("copying stub files...")
	for _, src := range s.stubs {
		rel := strings.TrimPrefix(src, filepath.Join(s.source, "stubs"))
		dest := s.target + rel
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return err
		}
		sfile, err := os.Open(src)
		if err != nil {
			return err
		}
		defer sfile.Close()
		dfile, err := os.Open(dest)
		if err != nil {
			return err
		}
		defer dfile.Close()
		if _, err := io.Copy(dfile, sfile); err != nil {
			return err
		}
	}
	return nil
}

func RunCommandsAction(sd *cli.SourdoughConfig, s *StarterConfig) error {
	if cgroup, ok := s.commands["default"]; ok {
		if err := RunCommandGroup(sd.Ctx, cgroup); err != nil {
			return err
		}
	}
	return nil
}
