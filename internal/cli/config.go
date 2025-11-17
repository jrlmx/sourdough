package cli

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/jrlmx/sourdough/internal/cleanup"
)

type Flags struct {
	Force    bool
	Preserve bool
}

type SourdoughConfig struct {
	Ctx   context.Context
	Cmd   string
	Args  []string
	Flags Flags
	CM    cleanup.Manager
}

func NewSourdoughConfig(ctx context.Context) *SourdoughConfig {
	var flags Flags
	flag.BoolVar(&flags.Force, "force", false, "force the operation")
	flag.BoolVar(&flags.Preserve, "preserve", false, "preserves artifacts if the operation fails")
	flag.Parse()
	parts := flag.Args()
	var command string
	var args []string
	if len(parts) > 0 {
		command = parts[0]
	}
	if len(parts) > 1 {
		args = parts[1:]
	}
	return &SourdoughConfig{
		Ctx:   ctx,
		Cmd:   command,
		Args:  args,
		Flags: flags,
		CM:    *cleanup.NewManager(),
	}
}

func (sc *SourdoughConfig) DataPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(homeDir, ".sourdough")
}

func (sc *SourdoughConfig) StarterPath() string {
	return filepath.Join(sc.DataPath(), "starters")
}

func (sc *SourdoughConfig) StarterOptions() []string {
	spath := sc.StarterPath()
	sdir, err := os.ReadDir(spath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatal("starter directory does not exist. try running 'sourdough install'")
		}
		log.Fatal(err)
	}
	var starters []string
	for _, entry := range sdir {
		if entry.IsDir() {
			starters = append(starters, entry.Name())
		}
	}
	return starters
}
