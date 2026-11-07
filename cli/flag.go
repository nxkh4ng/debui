package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nxkh4ng/debui/logger"
	"github.com/spf13/pflag"
)

type Options struct {
	Version bool
	Driver  string
	Port    int
	DSN     string
}

func New(version string) Options {
	fs := pflag.NewFlagSet("debui", pflag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var opts Options

	fs.BoolVarP(&opts.Version, "version", "v", false, "print version and exit")
	fs.StringVarP(&opts.Driver, "driver", "d", "sqlite", "database driver: sqlite | postgres")
	fs.IntVarP(&opts.Port, "port", "p", 8899, "port for the web UI: 1 < port < 65535")
	fs.StringVar(&opts.DSN, "dsn", "", "database connection string (required)")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags]\n", filepath.Base(os.Args[0]))
		fs.PrintDefaults()
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, pflag.ErrHelp) {
			os.Exit(0)
		}

		logger.Errorf("%v\n\n", err)
		fs.Usage()
		os.Exit(2)
	}

	if opts.Version {
		fmt.Printf("debui version %s\n", version)
		os.Exit(0)
	}

	if opts.DSN == "" {
		logger.Errorf("--dsn is required, currently empty\n")
		os.Exit(2)
	}

	switch opts.Driver {
	case "sqlite", "postgres":
	default:
		logger.Errorf("unsupported driver %q (want: sqlite, postgres)\n", opts.Driver)
		os.Exit(2)
	}

	if opts.Port < 1 || opts.Port > 65535 {
		logger.Errorf("port is out of range: %d (want: 1 < port < 65535)\n", opts.Port)
		os.Exit(2)
	}

	return opts
}
