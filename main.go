package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nxkh4ng/debui/driver"
	"github.com/nxkh4ng/debui/server"
	"github.com/spf13/pflag"
)

const (
	dbTimeout       = 15 * time.Second
	shutdownTimeout = 30 * time.Second
)

var VERSION = "dev"

type Options struct {
	Driver  string
	Port    int
	DSN     string
	Version bool
}

func main() {
	var opts Options

	pflag.StringVarP(&opts.Driver, "driver", "d", "sqlite", "database driver: sqlite | postgres")
	pflag.StringVar(&opts.DSN, "dsn", "", "database connection string (required)")
	pflag.IntVarP(&opts.Port, "port", "p", 8899, "port for the web UI: 1 < port < 65535")
	pflag.BoolVarP(&opts.Version, "version", "v", false, "print version and exit")
	pflag.Parse()

	if opts.Version {
		fmt.Printf("debui version %s\n", VERSION)
		os.Exit(0)
	}

	if opts.DSN == "" {
		fmt.Fprintf(os.Stderr, "ERROR: --dsn is required, currently empty\n\n")
		pflag.Usage()
		os.Exit(1)
	}

	switch opts.Driver {
	case "sqlite", "postgres":
	default:
		fmt.Fprintf(os.Stderr, "ERROR: unsupported drvier %q\n\n", opts.Driver)
		pflag.Usage()
		os.Exit(1)
	}

	if opts.Port < 1 || opts.Port > 65535 {
		fmt.Fprintf(os.Stderr, "ERROR: port is out of range: %d\n\n", opts.Port)
		pflag.Usage()
		os.Exit(1)
	}

	if code := run(opts); code != 0 {
		os.Exit(code)
	}
}

func run(opts Options) int {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dbCtx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	drv, err := driver.New(dbCtx, opts.Driver, opts.DSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		return 1
	}
	defer func() { _ = drv.Close() }()
	fmt.Printf("%s connected\n", opts.Driver)

	addr := fmt.Sprintf("127.0.0.1:%d", opts.Port)
	srv := server.New(addr, drv)

	fmt.Printf("server listening on http://%s\n", addr)
	if err := srv.RunAndShut(ctx, shutdownTimeout); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		return 1
	}
	fmt.Println("server shutdown gracefully")

	return 0
}
