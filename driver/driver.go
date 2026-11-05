package driver

import (
	"context"
)

type Driver interface {
	ListTables(ctx context.Context) ([]string, error)
	Close() error
}

func New(ctx context.Context, driverName, dsn string) (Driver, error) {
	switch driverName {
	case "postgres":
		return OpenPostgres(ctx, dsn)
	default:
		return OpenSQLite(ctx, dsn)
	}
}
