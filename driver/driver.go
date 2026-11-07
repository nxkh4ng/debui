package driver

import (
	"context"
	"errors"
	"strings"
)

var ErrTableNotFound = errors.New("table not found")

type Column struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type TableData struct {
	Columns []Column `json:"columns"`
	Rows    [][]any  `json:"rows"`
	Total   int64    `json:"total"`
}

type Driver interface {
	ListTables(ctx context.Context) ([]string, error)
	GetData(ctx context.Context, table string, limit, offset int) (TableData, error)
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

func quoteTable(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}
