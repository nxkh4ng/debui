package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nxkh4ng/debui/cli"
	"github.com/nxkh4ng/debui/driver"
	"github.com/nxkh4ng/debui/logger"
	"github.com/nxkh4ng/debui/server"
)

const (
	dbTimeout       = 15 * time.Second
	shutdownTimeout = 30 * time.Second
)

var VERSION = "dev"

func main() {
	opts := cli.New(VERSION)

	if code := run(opts); code != 0 {
		os.Exit(code)
	}
}

func run(opts cli.Options) int {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dbCtx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	drv, err := driver.New(dbCtx, opts.Driver, opts.DSN)
	if err != nil {
		logger.Errorf("%v\n", err)
		return 1
	}
	defer func() {
		if err := drv.Close(); err != nil {
			logger.Warnf("close db: %v\n", err)
		}
	}()
	logger.Infof("%s connected\n", opts.Driver)

	addr := fmt.Sprintf("127.0.0.1:%d", opts.Port)
	srv := server.New(addr, drv)

	logger.Infof("server listening on http://%s\n", addr)
	if err := srv.RunAndShut(ctx, shutdownTimeout); err != nil {
		logger.Errorf("%v\n", err)
		return 1
	}
	logger.Infof("server shutdown gracefully\n")

	return 0
}
